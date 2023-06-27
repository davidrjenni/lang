// Copyright (c) 2023 David Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compiler // import "davidrjenni.io/lang/compiler"

import (
	"fmt"
	"io"

	"davidrjenni.io/lang/ir"
)

func Compile(out io.Writer, n ir.Node) {
	c := &compiler{out: out}
	fmt.Fprint(out, macros)
	fmt.Fprint(out, main)
	c.compile(n)
	fmt.Fprint(out, epilogue)
	fmt.Fprint(out, data)
}

type compiler struct {
	out io.Writer
}

func (c *compiler) compile(n ir.Node) {
	switch n := n.(type) {
	case *ir.BinaryExpr:
		c.printf("%s %s, %s", op(n.Op, n.RHS.Type), rval(n.LHS), reg(n.RHS))
	case ir.Call:
		c.printf("%s", n)
	case ir.CJump:
		c.printf("%s %s", CJump, n)
	case ir.Jump:
		c.printf("%s %s", Jump, n)
	case ir.Label:
		fmt.Fprintf(c.out, "%s:\n", n)
	case *ir.Load:
		c.printf("%s %s, %s", mov(n.Dst.Type), rval(n.Src), reg(n.Dst))
	case ir.Seq:
		for _, s := range n {
			c.compile(s)
		}
	case *ir.UnaryExpr:
		c.printf("%s %s", op(n.Op, n.Reg.Type), reg(n.Reg))
	default:
		panic(fmt.Sprintf("unexpected type %T", n))
	}
}

func reg(r *ir.Reg) string {
	switch r.Type {
	case ir.BoolReg:
		if r.Second {
			return "%bl"
		}
		return "%al"
	case ir.F64Reg:
		if r.Second {
			return "%xmm1"
		}
		return "%xmm0"
	case ir.I64Reg:
		if r.Second {
			return "%rbx"
		}
		return "%rax"
	default:
		panic(fmt.Sprintf("unexpected type %d", r.Type))
	}
}

func rval(v ir.RVal) string {
	switch v := v.(type) {
	case ir.Bool:
		if v {
			return "$1"
		}
		return "$0"
	case ir.I64:
		return fmt.Sprintf("$%d", v)
	case *ir.Reg:
		return reg(v)
	default:
		panic(fmt.Sprintf("unexpected type %T", v))
	}
}

func (c *compiler) printf(f string, args ...interface{}) {
	fmt.Fprintf(c.out, "\t%s\n", fmt.Sprintf(f, args...))
}

const main = `
	.section .text
	.global main
main:
	pushq %rbp
	movq %rsp, %rbp
	subq $8, %rsp
`

const epilogue = `
	movq $0, %rax
	movq %rbp, %rsp
	popq %rbp
	ret
`

const macros = `
.macro AssertViolated
    movq $___fmt_assert, %rdi
    movq $0, %rax
    call printf
    movq $1, %rdi
    movq $0, %rax
    call exit
.endm
`

const data = `
	.section .data
___fmt_assert: .string "Assertion violated\n"
`
