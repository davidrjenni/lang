// Copyright (c) 2023 David Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ir // import "davidrjenni.io/lang/ir"

import (
	"fmt"
	"strconv"

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

func Translate(b *ast.Block, info types.Info, passes ...Pass) *Frame {
	t := &translator{
		info:   info,
		passes: passes,
	}
	t.translateFrame(b, Label("main"))

	return t.frames[0]
}

type translator struct {
	info   types.Info
	passes []Pass

	labels      int
	frameStates []*frameState

	frames []*Frame
}

// frameState represents the per-frame translator state.
type frameState struct {
	stack     int
	vars      map[string]int
	forStarts []Label
	forEnds   []Label
}

func (t *translator) translateFrame(b *ast.Block, label Label) {
	fs := &frameState{
		vars: make(map[string]int),
	}

	t.frameStates = append(t.frameStates, fs)

	s := t.translateCmd(b)
	s = flatten(s)
	for _, p := range t.passes {
		s = p(s)
	}

	i := len(t.frameStates) - 1
	t.frameStates = t.frameStates[:i]

	t.frames = append(t.frames, &Frame{Name: label, Seq: s, Stack: -fs.stack})
}

func (t *translator) translateCmd(cmd ast.Cmd) Seq {
	switch cmd := cmd.(type) {
	case *ast.Assert:
		return t.translateAssert(cmd)
	case *ast.Assign:
		return t.translateAssign(cmd)
	case *ast.Block:
		return t.translateBlock(cmd)
	case *ast.Break:
		return t.translateBreak(cmd)
	case *ast.Continue:
		return t.translateContinue(cmd)
	case *ast.For:
		return t.translateFor(cmd)
	case *ast.If:
		return t.translateIf(cmd)
	case *ast.VarDecl:
		return t.translateVarDecl(cmd)
	default:
		panic(fmt.Sprintf("unexpected type %T", cmd))
	}
}

func (t *translator) translateAssert(a *ast.Assert) Seq {
	label := t.label()
	return Seq{
		t.boolCheck(a.X, true_),
		&CJump{Label: label, pos: a.Pos()},
		&Load{Src: I64(a.Pos().Line), Dst: i64Reg2, pos: a.Pos()},
		&Call{Label: assertViolated, pos: a.Pos()},
		label,
	}
}

func (t *translator) translateAssign(a *ast.Assign) Seq {
	src := t.translateRVal(a.X)
	sz := t.info.Uses[a.Ident].Type.Size()
	off := t.fs().vars[a.Ident.Name]
	mem := &Mem{Off: off}
	store := &Store{Src: src, Dst: mem, Size: regType(sz), pos: a.Pos()}
	return Seq{store}
}

func (t *translator) translateBlock(b *ast.Block) (s Seq) {
	for _, c := range b.Cmds {
		n := t.translateCmd(c)
		s = append(s, n)
	}
	return s
}

func (t *translator) translateBreak(b *ast.Break) Seq {
	return Seq{
		&Jump{Label: t.fs().forEnds[len(t.fs().forEnds)-1], pos: b.Pos()},
	}
}

func (t *translator) translateContinue(c *ast.Continue) Seq {
	return Seq{
		&Jump{Label: t.fs().forStarts[len(t.fs().forStarts)-1], pos: c.Pos()},
	}
}

func (t *translator) translateFor(f *ast.For) Seq {
	start := t.label()
	end := t.label()
	t.fs().forStarts = append(t.fs().forStarts, start)
	t.fs().forEnds = append(t.fs().forEnds, end)
	seq := Seq{
		start,
		t.boolCheck(f.X, false_),
		&CJump{Label: end, pos: f.Pos()},
		t.translateCmd(f.Block),
		&Jump{Label: start, pos: f.Pos()},
		end,
	}
	i := len(t.fs().forStarts) - 1
	t.fs().forStarts = t.fs().forStarts[:i]
	t.fs().forEnds = t.fs().forEnds[:i]
	return seq
}

func (t *translator) translateIf(i *ast.If) Seq {
	end := t.label()
	seq := Seq{
		t.boolCheck(i.X, false_),
		&CJump{Label: end, pos: i.Pos()},
		t.translateCmd(i.Block),
	}
	if i.Else != nil {
		endElse := t.label()
		seq = append(seq, Seq{
			&Jump{Label: endElse, pos: i.Pos()},
			end,
			t.translateCmd(i.Else.Cmd),
		})
		end = endElse
	}
	return append(seq, end)
}

func (t *translator) translateVarDecl(d *ast.VarDecl) Seq {
	src := t.translateRVal(d.X)
	sz := t.info.Uses[d.Ident].Type.Size()
	t.fs().stack -= sz
	t.fs().vars[d.Ident.Name] = t.fs().stack
	mem := &Mem{Off: t.fs().stack}
	store := &Store{Src: src, Dst: mem, Size: regType(sz), pos: d.Pos()}
	return Seq{store}
}

func (t *translator) boolCheck(x ast.Expr, expect Bool) Seq {
	check := t.translateRVal(x)
	return Seq{
		&Load{Src: check, Dst: boolReg1, pos: x.Pos()},
		&BinaryInstr{RHS: boolReg1, Op: Cmp, LHS: expect, pos: x.Pos()},
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
			seq = append(seq, &UnaryInstr{Reg: pr, Op: Push, pos: x.Pos()})
		}

		// Load LHS into the first register.
		lhs := t.translateRVal(x.LHS)
		seq = append(seq, &Load{Src: lhs, Dst: r1, pos: x.Pos()})

		// Load RHS into the second register.
		if pushed {
			seq = append(seq, &UnaryInstr{Reg: i64Reg2, Op: Pop, pos: x.Pos()})
		} else {
			seq = append(seq, &Load{Src: rhs, Dst: r2, pos: x.Pos()})
		}

		seq = append(seq, &BinaryInstr{RHS: r1, Op: binOp(x.Op), LHS: r2, pos: x.Pos()})
		if isCmp(x.Op) {
			r1 = boolReg1
			seq = append(seq, &UnaryInstr{Reg: r1, Op: cmpOp(x.Op), pos: x.Pos()})
		}
		return &seqExpr{Seq: seq, Dst: r1}
	case *ast.Bool:
		return Bool(x.Val == "true")
	case *ast.F64:
		val, err := strconv.ParseFloat(x.Val, 64)
		if err != nil {
			panic(fmt.Sprintf("cannot convert f64: %v", err))
		}
		return F64(val)
	case *ast.I64:
		val, err := strconv.ParseInt(x.Val, 10, 0)
		if err != nil {
			panic(fmt.Sprintf("cannot convert i64: %v", err))
		}
		return I64(val)
	case *ast.Ident:
		return &Mem{Off: t.fs().vars[x.Name]}
	case *ast.ParenExpr:
		return t.translateRVal(x.X)
	case *ast.UnaryExpr:
		switch x.Op {
		case lexer.Minus:
			val := t.translateRVal(x.X)
			return &seqExpr{
				Seq: Seq{
					&Load{Src: val, Dst: i64Reg1, pos: x.Pos()},
					&UnaryInstr{Op: Neg, Reg: i64Reg1, pos: x.Pos()},
				},
				Dst: i64Reg1,
			}
		case lexer.Not:
			return &seqExpr{
				Seq: Seq{
					t.boolCheck(x.X, true_),
					&UnaryInstr{Op: Setne, Reg: boolReg1, pos: x.Pos()},
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

func (t *translator) fs() *frameState {
	return t.frameStates[len(t.frameStates)-1]
}
