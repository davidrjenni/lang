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
	intReg1  = &Reg{Type: I64Reg}
	boolReg1 = &Reg{Type: BoolReg}
)

func Translate(b *ast.Block) Node {
	t := &translator{}
	return t.translateCmd(b)
}

type translator struct {
	labels int
}

func (t *translator) translateCmd(cmd ast.Cmd) Node {
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
	return Seq{
		&Load{
			Src: t.translateRVal(x),
			Dst: boolReg1,
		},
		&BinaryExpr{
			RHS: boolReg1,
			Op:  Cmp,
			LHS: expect,
		},
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
			return &SeqExpr{
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
