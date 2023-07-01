// Copyright (c) 2023 David Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ast // import "davidrjenni.io/lang/ast"

import (
	"fmt"
	"io"
	"strings"
)

func Dump(out io.Writer, n Node) {
	d := &dumper{out: out, indent: 0}
	d.dump(n)
}

type dumper struct {
	out    io.Writer
	indent int
}

func (d *dumper) dump(n Node) {
	switch n := n.(type) {
	case Cmd:
		d.dumpCmd(n)
	case Expr:
		d.dumpExpr(n)
	default:
		panic(fmt.Sprintf("unexpected type %T", n))
	}
}

func (d *dumper) dumpCmd(cmd Cmd) {
	switch cmd := cmd.(type) {
	case *Assert:
		d.enter("Assert(")
		d.dumpPos(cmd)
		d.print("X: ")
		d.dumpExpr(cmd.X)
		d.exit(")")
	case *Block:
		d.enter("Block(")
		d.dumpPos(cmd)
		for i, c := range cmd.Cmds {
			d.printf("%d: ", i)
			d.dumpCmd(c)
			d.println()
		}
		d.exit(")")
	case *For:
		d.enter("For(")
		d.dumpPos(cmd)
		d.print("X: ")
		d.dumpExpr(cmd.X)
		d.println()
		d.print("Block: ")
		d.dumpCmd(cmd.Block)
		d.exit(")")
	case *If:
		d.enter("If(")
		d.dumpPos(cmd)
		d.print("X: ")
		d.dumpExpr(cmd.X)
		d.println()
		d.print("Block: ")
		d.dumpCmd(cmd.Block)
		d.exit(")")
	default:
		panic(fmt.Sprintf("unexpected type %T", cmd))
	}
}

func (d *dumper) dumpExpr(x Expr) {
	switch x := x.(type) {
	case *BinaryExpr:
		d.enter("BinaryExpr(")
		d.dumpPos(x)
		d.print("LHS: ")
		d.dump(x.LHS)
		d.println()
		d.printf("Op: %s", x.Op.String())
		d.println()
		d.print("RHS: ")
		d.dump(x.RHS)
		d.exit(")")
	case Lit:
		d.dumpLit(x)
	case *ParenExpr:
		d.enter("ParenExpr(")
		d.dumpPos(x)
		d.print("X: ")
		d.dump(x.X)
		d.exit(")")
	case *UnaryExpr:
		d.enter("UnaryExpr(")
		d.dumpPos(x)
		d.printf("Op: %s", x.Op.String())
		d.println()
		d.print("X: ")
		d.dump(x.X)
		d.exit(")")
	default:
		panic(fmt.Sprintf("unexpected type %T", x))
	}
}

func (d *dumper) dumpLit(l Lit) {
	switch l := l.(type) {
	case *Bool:
		d.printf("Bool(Val: %v, Pos: %s, End: %s)", l.Val, l.Pos(), l.End())
	case *F64:
		d.printf("F64(Val: %v, Pos: %s, End: %s)", l.Val, l.Pos(), l.End())
	case *I64:
		d.printf("I64(Val: %v, Pos: %s, End: %s)", l.Val, l.Pos(), l.End())
	case *String:
		d.printf("String(Val: %q, Pos: %s, End: %s)", l.Val, l.Pos(), l.End())
	default:
		panic(fmt.Sprintf("unexpected type %T", l))
	}
}

func (d *dumper) dumpPos(n Node) {
	d.printf("Pos: (Start: %s, End: %s)", n.Pos(), n.End())
	d.println()
}

func (d *dumper) enter(s string) {
	d.print(s)
	d.indent++
	d.println()
}

func (d *dumper) exit(s string) {
	d.indent--
	d.println()
	d.print(s)
}

func (d *dumper) print(s string) { fmt.Fprint(d.out, s) }

func (d *dumper) println() { fmt.Fprint(d.out, "\n"+strings.Repeat("\t", d.indent)) }

func (d *dumper) printf(f string, args ...interface{}) {
	fmt.Fprintf(d.out, f, args...)
}
