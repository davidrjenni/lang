// Copyright (c) 2023 David Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ir // import "davidrjenni.io/lang/ir"

import (
	"fmt"

	"davidrjenni.io/lang/lexer"
)

//go:generate stringer -type=Op -linecomment
//go:generate stringer -type=RegType -linecomment

type (
	Node interface {
		node()
	}

	Frame struct {
		Name  Label
		Seq   Seq
		Stack int
	}

	Label string

	Seq []Node
)

func (*Frame) node() {}
func (Label) node()  {}
func (Seq) node()    {}

type (
	Cmd interface {
		Pos() lexer.Pos
		cmd()
		Node
	}

	BinaryInstr struct {
		RHS *Reg
		Op  Op
		LHS RVal
		pos lexer.Pos
	}

	Call struct {
		Label Label
		pos   lexer.Pos
	}

	CJump struct {
		Label Label
		pos   lexer.Pos
	}

	Jump struct {
		Label Label
		pos   lexer.Pos
	}

	Load struct {
		Src RVal
		Dst *Reg
		pos lexer.Pos
	}

	Store struct {
		Src RVal
		Dst *Mem
		pos lexer.Pos
	}

	UnaryInstr struct {
		Reg *Reg
		Op  Op
		pos lexer.Pos
	}
)

func (c *BinaryInstr) Pos() lexer.Pos { return c.pos }
func (c *Call) Pos() lexer.Pos        { return c.pos }
func (c *CJump) Pos() lexer.Pos       { return c.pos }
func (c *Jump) Pos() lexer.Pos        { return c.pos }
func (c *Load) Pos() lexer.Pos        { return c.pos }
func (c *Store) Pos() lexer.Pos       { return c.pos }
func (c *UnaryInstr) Pos() lexer.Pos  { return c.pos }

func (*BinaryInstr) node() {}
func (*Call) node()        {}
func (*CJump) node()       {}
func (*Jump) node()        {}
func (*Load) node()        {}
func (*Store) node()       {}
func (*UnaryInstr) node()  {}

func (*BinaryInstr) cmd() {}
func (*Call) cmd()        {}
func (*CJump) cmd()       {}
func (*Jump) cmd()        {}
func (*Load) cmd()        {}
func (*Store) cmd()       {}
func (*UnaryInstr) cmd()  {}

type (
	RVal interface {
		rval()
		Node
	}

	Bool bool

	F64 float64

	I64 int64

	seqExpr struct {
		Seq Seq
		Dst *Reg
	}
)

func (Bool) node()     {}
func (F64) node()      {}
func (I64) node()      {}
func (*seqExpr) node() {}

func (Bool) rval()     {}
func (F64) rval()      {}
func (I64) rval()      {}
func (*seqExpr) rval() {}

type (
	LVal interface {
		lval()
		RVal
	}

	Mem struct {
		Off int
	}

	Reg struct {
		Type   RegType
		Second bool
	}
)

func (*Mem) node() {}
func (*Reg) node() {}

func (*Mem) rval() {}
func (*Reg) rval() {}

func (*Mem) lval() {}
func (*Reg) lval() {}

type Op int

const (
	Push Op = iota // push
	Pop            // pop

	Neg // neg

	Add // add
	Sub // sub
	Mul // mul
	Div // div

	Cmp // cmp
	And // and
	Or  // or

	Setl  // setl
	Setle // setle
	Sete  // sete
	Setne // setne
	Setg  // setg
	Setge // setge
)

var (
	binOps = map[lexer.Tok]Op{
		lexer.Plus:     Add,
		lexer.Minus:    Sub,
		lexer.Multiply: Mul,
		lexer.Divide:   Div,

		lexer.And:     And,
		lexer.Or:      Or,
		lexer.Implies: Or, // Assumes desugaring of a => b into ¬a ∨ b.

		lexer.Less:      Cmp,
		lexer.LessEq:    Cmp,
		lexer.Equal:     Cmp,
		lexer.NotEqual:  Cmp,
		lexer.Greater:   Cmp,
		lexer.GreaterEq: Cmp,
	}

	cmpOps = map[lexer.Tok]Op{
		lexer.Less:      Setl,
		lexer.LessEq:    Setle,
		lexer.Equal:     Sete,
		lexer.NotEqual:  Setne,
		lexer.Greater:   Setg,
		lexer.GreaterEq: Setge,
	}
)

func binOp(op lexer.Tok) Op {
	o, ok := binOps[op]
	if !ok {
		panic(fmt.Sprintf("unexpected op %s", op))
	}
	return o
}

func cmpOp(op lexer.Tok) Op {
	o, ok := cmpOps[op]
	if !ok {
		panic(fmt.Sprintf("unexpected op %s", op))
	}
	return o
}

func isCmp(op lexer.Tok) bool {
	return lexer.Less <= op && op <= lexer.GreaterEq
}

type RegType int

const (
	BoolReg RegType = iota // bool
	F64Reg                 // f64
	I64Reg                 // i64
)
