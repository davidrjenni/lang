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
	case *Comment:
		d.printf("Comment(Text: %q, Pos: %s, End: %s)", n.Text, n.Pos(), n.End())
	case Expr:
		d.dumpExpr(n)
	default:
		panic(fmt.Sprintf("unexpected type %T", n))
	}
}

func (d *dumper) dumpTypes(ts []Type) {
	for i, t := range ts {
		d.printf("%d: ", i)
		d.dumpType(t)
		d.println()
	}
}

func (d *dumper) dumpType(t Type) {
	switch t := t.(type) {
	case *Func:
		d.enter("Func(")
		d.dumpPos(t)
		d.enter("Params: (")
		d.dumpTypes(t.Params)
		d.exit(")")
		d.println()
		d.print("Result: ")
		d.dumpType(t.Result)
		d.exit(")")
	case *Scalar:
		d.enter("Scalar(")
		d.dumpPos(t)
		d.printf("Name: %s", t.Name)
		d.exit(")")
	default:
		panic(fmt.Sprintf("unexpected type %T", t))
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
	case *Assign:
		d.enter("Assign(")
		d.dumpPos(cmd)
		d.print("Ident: ")
		d.dumpExpr(cmd.Ident)
		d.println()
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
	case *Break:
		d.printf("Break(Pos: %s, End: %s)", cmd.Pos(), cmd.End())
	case *Continue:
		d.printf("Continue(Pos: %s, End: %s)", cmd.Pos(), cmd.End())
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
		if cmd.Else != nil {
			d.println()
			d.print("Else: ")
			d.enter("Else(")
			d.dumpPos(cmd.Else)
			d.print("Cmd: ")
			d.dumpCmd(cmd.Else.Cmd)
			d.exit(")")
		}
		d.exit(")")
	case *Return:
		d.enter("Return(")
		d.dumpPos(cmd)
		d.print("X: ")
		d.dumpExpr(cmd.X)
		d.exit(")")
	case *VarDecl:
		d.enter("Var(")
		d.dumpPos(cmd)
		d.print("Ident: ")
		d.dumpExpr(cmd.Ident)
		d.println()
		d.print("X: ")
		d.dumpExpr(cmd.X)
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
	case *Ident:
		d.printf("Ident(Name: %q, Pos: %s, End: %s)", x.Name, x.Pos(), x.End())
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
	case *FuncLit:
		d.enter("FuncLit(")
		d.dumpPos(l)
		d.enter("Params: (")
		d.dumpFields(l.Params)
		d.exit(")")
		d.println()
		d.print("Result: ")
		d.dumpType(l.Result)
		d.println()
		d.dumpCmd(l.Block)
		d.exit(")")
	case *I64:
		d.printf("I64(Val: %v, Pos: %s, End: %s)", l.Val, l.Pos(), l.End())
	case *String:
		d.printf("String(Val: %q, Pos: %s, End: %s)", l.Val, l.Pos(), l.End())
	default:
		panic(fmt.Sprintf("unexpected type %T", l))
	}
}

func (d *dumper) dumpFields(fs []*Field) {
	for i, f := range fs {
		d.printf("%d: ", i)
		d.enter("Field(")
		d.dumpPos(f)
		d.print("Ident: ")
		d.dumpExpr(f.Ident)
		d.println()
		d.print("Type: ")
		d.dumpType(f.Type)
		d.exit(")")
		d.println()
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
