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

	Jump  // jmp
	CJump // je

	Cmpq // cmpq
	Cmpb // cmpb

	Setne // setne

	Call // call
)

var ops = map[ir.Op]map[ir.RegType]Op{
	ir.Cmp:   {ir.I64Reg: Cmpq, ir.BoolReg: Cmpb},
	ir.Setne: {ir.BoolReg: Setne},
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
