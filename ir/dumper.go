// Copyright (c) 2023 David Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ir // import "davidrjenni.io/lang/ir"

import (
	"fmt"
	"io"
)

func Dump(out io.Writer, n Node) {
	d := &dumper{out: out}
	d.dump(n)
}

type dumper struct {
	out io.Writer
}

func (d *dumper) dump(n Node) {
	switch n := n.(type) {
	case *BinaryExpr:
		seqx, ok := n.LHS.(*SeqExpr)
		lhs := n.LHS
		if ok {
			d.dump(seqx.Seq)
			lhs = seqx.Dst
		}
		d.printf("%s %s %s", n.Op, lval(n.RHS), rval(lhs))
	case Call:
		d.printf("call %s", n)
	case CJump:
		d.printf("cjump %s", n)
	case Jump:
		d.printf("jump %s", n)
	case Label:
		d.printf("%s", n)
	case *Load:
		seqx, ok := n.Src.(*SeqExpr)
		src := n.Src
		if ok {
			d.dump(seqx.Seq)
			src = seqx.Dst
		}
		d.printf("load %s <- %s", lval(n.Dst), rval(src))
	case Seq:
		for _, s := range n {
			d.dump(s)
		}
	case *UnaryExpr:
		d.printf("%s %s", n.Op, lval(n.Reg))
	default:
		panic(fmt.Sprintf("unexpected type %T", n))
	}
}

func rval(n RVal) string {
	switch n := n.(type) {
	case Bool:
		return fmt.Sprintf("bool(%v)", n)
	case F64:
		return fmt.Sprintf("f64(%v)", n)
	case I64:
		return fmt.Sprintf("i64(%d)", n)
	case LVal:
		return lval(n)
	default:
		panic(fmt.Sprintf("unexpected type %T", n))
	}
}

func lval(n LVal) string {
	switch n := n.(type) {
	case *Mem:
		return fmt.Sprintf("m%d", n.Off)
	case *Reg:
		i := 0
		if n.Second {
			i = 1
		}
		return fmt.Sprintf("r%s.%d", n.Type, i)
	default:
		panic(fmt.Sprintf("unexpected type %T", n))
	}
}

func (d *dumper) printf(f string, args ...interface{}) {
	fmt.Fprintf(d.out, "%s\n", fmt.Sprintf(f, args...))
}