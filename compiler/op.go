// Copyright (c) 2023 David Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compiler // import "davidrjenni.io/lang/compiler"

import (
	"fmt"

	"davidrjenni.io/lang/ir"
)

//go:generate stringer -type=Op -linecomment

type Op int

const (
	Movq Op = iota // movq
	Movb           // movb

	Push // pushq
	Pop  // popq

	Jump  // jmp
	CJump // je

	Neg // negq

	Add // addq
	Sub // subq
	Mul // imulq
	Div // idiv

	And // andb
	Or  // orb

	Cmpq // cmpq
	Cmpb // cmpb

	Setl  // setl
	Setle // setle
	Sete  // sete
	Setne // setne
	Setg  // setg
	Setge // setge

	Call // call
)

var ops = map[ir.Op]map[ir.RegType]Op{
	ir.Push: {ir.I64Reg: Push},
	ir.Pop:  {ir.I64Reg: Pop},

	ir.Neg: {ir.I64Reg: Neg},

	ir.Add: {ir.I64Reg: Add},
	ir.Sub: {ir.I64Reg: Sub},
	ir.Mul: {ir.I64Reg: Mul},
	ir.Div: {ir.I64Reg: Div},

	ir.And: {ir.BoolReg: And},
	ir.Or:  {ir.BoolReg: Or},
	ir.Cmp: {ir.I64Reg: Cmpq, ir.BoolReg: Cmpb},

	ir.Setl:  {ir.BoolReg: Setl},
	ir.Setle: {ir.BoolReg: Setle},
	ir.Sete:  {ir.BoolReg: Sete},
	ir.Setne: {ir.BoolReg: Setne},
	ir.Setg:  {ir.BoolReg: Setg},
	ir.Setge: {ir.BoolReg: Setge},
}

func mov(t ir.RegType) Op {
	switch t {
	case ir.BoolReg:
		return Movb
	case ir.I64Reg:
		return Movq
	default:
		panic(fmt.Sprintf("unexpected reg %s", t))
	}
}

func op(op ir.Op, t ir.RegType) Op {
	candidates, ok := ops[op]
	if !ok {
		panic(fmt.Sprintf("unexpected op %s", op))
	}
	o, ok := candidates[t]
	if !ok {
		panic(fmt.Sprintf("unexpected %s reg for op %s", t, op))
	}
	return o
}
