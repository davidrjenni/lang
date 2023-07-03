// Copyright (c) 2023 David Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ast // import "davidrjenni.io/lang/ast"

import "davidrjenni.io/lang/lexer"

type Node interface {
	Pos() lexer.Pos
	End() lexer.Pos
	node()
}

type Comment struct {
	Text     string
	StartPos lexer.Pos
	EndPos   lexer.Pos
}

func (c *Comment) Pos() lexer.Pos { return c.StartPos }
func (c *Comment) End() lexer.Pos { return c.EndPos }

func (*Comment) node() {}

type Else struct {
	Cmd      Cmd
	StartPos lexer.Pos
}

func (e *Else) Pos() lexer.Pos { return e.StartPos }
func (e *Else) End() lexer.Pos { return e.Cmd.End() }

func (*Else) node() {}

type (
	Cmd interface {
		cmd()
		Node
	}

	Assert struct {
		X        Expr
		StartPos lexer.Pos
		EndPos   lexer.Pos
	}

	Block struct {
		Cmds     []Cmd
		StartPos lexer.Pos
		EndPos   lexer.Pos
	}

	Break struct {
		StartPos lexer.Pos
		EndPos   lexer.Pos
	}

	Continue struct {
		StartPos lexer.Pos
		EndPos   lexer.Pos
	}

	For struct {
		X        Expr
		Block    *Block
		StartPos lexer.Pos
	}

	If struct {
		X        Expr
		Block    *Block
		Else     *Else
		StartPos lexer.Pos
	}
)

func (c *Assert) Pos() lexer.Pos { return c.StartPos }
func (c *Assert) End() lexer.Pos { return c.EndPos }

func (c *Block) Pos() lexer.Pos { return c.StartPos }
func (c *Block) End() lexer.Pos { return c.EndPos }

func (c *Break) Pos() lexer.Pos { return c.StartPos }
func (c *Break) End() lexer.Pos { return c.EndPos }

func (c *Continue) Pos() lexer.Pos { return c.StartPos }
func (c *Continue) End() lexer.Pos { return c.EndPos }

func (c *For) Pos() lexer.Pos { return c.StartPos }
func (c *For) End() lexer.Pos { return c.Block.End() }

func (c *If) Pos() lexer.Pos { return c.StartPos }
func (c *If) End() lexer.Pos {
	if c.Else != nil {
		return c.Else.End()
	}
	return c.Block.End()
}

func (*Assert) node()   {}
func (*Block) node()    {}
func (*Break) node()    {}
func (*Continue) node() {}
func (*For) node()      {}
func (*If) node()       {}

func (*Assert) cmd()   {}
func (*Block) cmd()    {}
func (*Break) cmd()    {}
func (*Continue) cmd() {}
func (*For) cmd()      {}
func (*If) cmd()       {}

type (
	Expr interface {
		expr()
		Node
	}

	BinaryExpr struct {
		LHS Expr
		Op  lexer.Tok
		RHS Expr
	}

	Ident struct {
		Name     string
		StartPos lexer.Pos
	}

	ParenExpr struct {
		X        Expr
		StartPos lexer.Pos
		EndPos   lexer.Pos
	}

	UnaryExpr struct {
		Op       lexer.Tok
		X        Expr
		StartPos lexer.Pos
	}
)

func (x *BinaryExpr) Pos() lexer.Pos { return x.LHS.Pos() }
func (x *BinaryExpr) End() lexer.Pos { return x.RHS.End() }

func (x *Ident) Pos() lexer.Pos { return x.StartPos }
func (x *Ident) End() lexer.Pos { return x.StartPos.Shift(len(x.Name)) }

func (x *ParenExpr) Pos() lexer.Pos { return x.StartPos }
func (x *ParenExpr) End() lexer.Pos { return x.EndPos }

func (x *UnaryExpr) Pos() lexer.Pos { return x.StartPos }
func (x *UnaryExpr) End() lexer.Pos { return x.X.End() }

func (*BinaryExpr) node() {}
func (*Ident) node()      {}
func (*ParenExpr) node()  {}
func (*UnaryExpr) node()  {}

func (*BinaryExpr) expr() {}
func (*Ident) expr()      {}
func (*ParenExpr) expr()  {}
func (*UnaryExpr) expr()  {}

type (
	Lit interface {
		lit()
		Expr
	}

	Bool struct {
		Val      bool
		StartPos lexer.Pos
		EndPos   lexer.Pos
	}

	F64 struct {
		Val      float64
		StartPos lexer.Pos
		EndPos   lexer.Pos
	}

	I64 struct {
		Val      int64
		StartPos lexer.Pos
		EndPos   lexer.Pos
	}

	String struct {
		Val      string
		StartPos lexer.Pos
		EndPos   lexer.Pos
	}
)

func (l *Bool) Pos() lexer.Pos { return l.StartPos }
func (l *Bool) End() lexer.Pos { return l.EndPos }

func (l *F64) Pos() lexer.Pos { return l.StartPos }
func (l *F64) End() lexer.Pos { return l.EndPos }

func (l *I64) Pos() lexer.Pos { return l.StartPos }
func (l *I64) End() lexer.Pos { return l.EndPos }

func (l *String) Pos() lexer.Pos { return l.StartPos }
func (l *String) End() lexer.Pos { return l.EndPos }

func (*Bool) node()   {}
func (*F64) node()    {}
func (*I64) node()    {}
func (*String) node() {}

func (*Bool) expr()   {}
func (*F64) expr()    {}
func (*I64) expr()    {}
func (*String) expr() {}

func (*Bool) lit()   {}
func (*F64) lit()    {}
func (*I64) lit()    {}
func (*String) lit() {}
