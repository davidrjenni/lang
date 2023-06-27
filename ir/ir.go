// Copyright (c) 2023 David Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ir // import "davidrjenni.io/lang/ir"

//go:generate stringer -type=Op -linecomment
//go:generate stringer -type=RegType -linecomment

type (
	Node interface {
		node()
	}

	BinaryExpr struct {
		RHS *Reg
		Op  Op
		LHS RVal
	}

	Call Label

	CJump Label

	Jump Label

	Label string

	Load struct {
		Src RVal
		Dst *Reg
	}

	Seq []Node

	Store struct {
		Src RVal
		Dst *Mem
	}

	UnaryExpr struct {
		Reg *Reg
		Op  Op
	}
)

func (*BinaryExpr) node() {}
func (Call) node()        {}
func (CJump) node()       {}
func (Jump) node()        {}
func (Label) node()       {}
func (*Load) node()       {}
func (Seq) node()         {}
func (*Store) node()      {}
func (*UnaryExpr) node()  {}

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

type RegType int

const (
	BoolReg RegType = iota // bool
	F64Reg                 // f64
	I64Reg                 // i64
)
