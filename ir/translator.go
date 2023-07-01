// Copyright (c) 2023 David Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ir // import "davidrjenni.io/lang/ir"

import (
	"fmt"

	"davidrjenni.io/lang/ast"
	"davidrjenni.io/lang/lexer"
	"davidrjenni.io/lang/types"
)

const (
	assertViolated = Label("AssertViolated")
)

const (
	true_  = Bool(true)
	false_ = Bool(false)
)

var (
	i64Reg1  = &Reg{Type: I64Reg, Second: false}
	i64Reg2  = &Reg{Type: I64Reg, Second: true}
	boolReg1 = &Reg{Type: BoolReg, Second: false}
	boolReg2 = &Reg{Type: BoolReg, Second: true}
	f64Reg1  = &Reg{Type: F64Reg, Second: false}
	f64Reg2  = &Reg{Type: F64Reg, Second: true}
)

func Translate(b *ast.Block, info types.Info, passes ...Pass) Seq {
	t := &translator{info: info}
	s := t.translateCmd(b)
	s = flatten(s)
	for _, p := range passes {
		s = p(s)
	}
	return s
}

type translator struct {
	info   types.Info
	labels int
}

func (t *translator) translateCmd(cmd ast.Cmd) Seq {
	switch cmd := cmd.(type) {
	case *ast.Assert:
		label := t.label()
		return Seq{
			t.boolCheck(cmd.X, true_),
			&CJump{Label: label, pos: cmd.Pos()},
			&Load{Src: I64(cmd.Pos().Line), Dst: i64Reg2, pos: cmd.Pos()},
			&Call{Label: assertViolated, pos: cmd.Pos()},
			label,
		}
	case *ast.Block:
		var seq Seq
		for _, c := range cmd.Cmds {
			n := t.translateCmd(c)
			seq = append(seq, n)
		}
		return seq
	case *ast.For:
		start := t.label()
		end := t.label()
		return Seq{
			start,
			t.boolCheck(cmd.X, false_),
			&CJump{Label: end, pos: cmd.Pos()},
			t.translateCmd(cmd.Block),
			&Jump{Label: start, pos: cmd.Pos()},
			end,
		}
	default:
		panic(fmt.Sprintf("unexpected type %T", cmd))
	}
}

func (t *translator) boolCheck(x ast.Expr, expect Bool) Seq {
	check := t.translateRVal(x)
	return Seq{
		&Load{Src: check, Dst: boolReg1, pos: x.Pos()},
		&BinaryExpr{RHS: boolReg1, Op: Cmp, LHS: expect, pos: x.Pos()},
	}
}

func (t *translator) translateRVal(x ast.Expr) RVal {
	switch x := x.(type) {
	case *ast.BinaryExpr:
		r1, r2 := i64Reg1, i64Reg2
		if types.Equal(t.info.Types[x.LHS].Type, &types.Bool{}) {
			r1, r2 = boolReg1, boolReg2
		}

		var seq Seq
		rhs := t.translateRVal(x.RHS)
		if seqx, ok := rhs.(*seqExpr); ok {
			seq = append(seq, seqx.Seq...)
			rhs = seqx.Dst
		}

		// Desugaring of a => b into ¬a ∨ b.
		if x.Op == lexer.Implies {
			x.LHS = &ast.UnaryExpr{X: x.LHS, Op: lexer.Not, StartPos: x.Pos()}
		}

		// Push RHS onto the stack.
		r, pushed := rhs.(*Reg)
		if pushed {
			pr := i64Reg1
			if r.Second {
				pr = i64Reg2
			}
			seq = append(seq, &UnaryExpr{Reg: pr, Op: Push, pos: x.Pos()})
		}

		// Load LHS into the first register.
		lhs := t.translateRVal(x.LHS)
		seq = append(seq, &Load{Src: lhs, Dst: r1, pos: x.Pos()})

		// Load RHS into the second register.
		if pushed {
			seq = append(seq, &UnaryExpr{Reg: i64Reg2, Op: Pop, pos: x.Pos()})
		} else {
			seq = append(seq, &Load{Src: rhs, Dst: r2, pos: x.Pos()})
		}

		seq = append(seq, &BinaryExpr{RHS: r1, Op: binOp(x.Op), LHS: r2, pos: x.Pos()})
		if isCmp(x.Op) {
			r1 = boolReg1
			seq = append(seq, &UnaryExpr{Reg: r1, Op: cmpOp(x.Op), pos: x.Pos()})
		}
		return &seqExpr{Seq: seq, Dst: r1}
	case *ast.Bool:
		return Bool(x.Val)
	case *ast.F64:
		return F64(x.Val)
	case *ast.I64:
		return I64(x.Val)
	case *ast.ParenExpr:
		return t.translateRVal(x.X)
	case *ast.UnaryExpr:
		switch x.Op {
		case lexer.Minus:
			val := t.translateRVal(x.X)
			return &seqExpr{
				Seq: Seq{
					&Load{Src: val, Dst: i64Reg1, pos: x.Pos()},
					&UnaryExpr{Op: Neg, Reg: i64Reg1, pos: x.Pos()},
				},
				Dst: i64Reg1,
			}
		case lexer.Not:
			return &seqExpr{
				Seq: Seq{
					t.boolCheck(x.X, true_),
					&UnaryExpr{Op: Setne, Reg: boolReg1, pos: x.Pos()},
				},
				Dst: boolReg1,
			}
		default:
			panic(fmt.Sprintf("unexpected operator %s", x.Op))
		}
	default:
		panic(fmt.Sprintf("unexpected type %T", x))
	}
}

func (t *translator) label() Label {
	t.labels++
	return Label(fmt.Sprintf(".L%d", t.labels))
}
