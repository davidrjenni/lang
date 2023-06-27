// Copyright (c) 2023 David Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ir // import "davidrjenni.io/lang/ir"

import (
	"fmt"

	"davidrjenni.io/lang/ast"
	"davidrjenni.io/lang/lexer"
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

func Translate(b *ast.Block, passes ...Pass) Seq {
	t := &translator{}
	s := t.translateCmd(b)
	s = flatten(s)
	for _, p := range passes {
		s = p(s)
	}
	return s
}

type translator struct {
	labels int
}

func (t *translator) translateCmd(cmd ast.Cmd) Seq {
	switch cmd := cmd.(type) {
	case *ast.Assert:
		label := t.label()
		return Seq{
			t.boolCheck(cmd.X, true_),
			CJump(label),
			Call(assertViolated),
			label,
		}
	case *ast.Block:
		var seq Seq
		for _, c := range cmd.Cmds {
			n := t.translateCmd(c)
			seq = append(seq, n)
		}
		return seq
	default:
		panic(fmt.Sprintf("unexpected type %T", cmd))
	}
}

func (t *translator) boolCheck(x ast.Expr, expect Bool) Seq {
	check := t.translateRVal(x)
	return Seq{
		&Load{Src: check, Dst: boolReg1},
		&BinaryExpr{RHS: boolReg1, Op: Cmp, LHS: expect},
	}
}

func (t *translator) translateRVal(x ast.Expr) RVal {
	switch x := x.(type) {
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
		case lexer.Not:
			return &seqExpr{
				Seq: Seq{
					t.boolCheck(x.X, true_),
					&UnaryExpr{Op: Setne, Reg: boolReg1},
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
