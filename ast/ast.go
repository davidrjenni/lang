// Copyright (c) 2023 David Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ast // import "davidrjenni.io/lang/ast"

import "davidrjenni.io/lang/lexer"

type Node interface {
	node()
}

type (
	Cmd interface {
		cmd()
		Node
	}

	Assert struct {
		X Expr
	}

	Block struct {
		Cmds []Cmd
	}
)

func (*Assert) node() {}
func (*Block) node()  {}

func (*Assert) cmd() {}
func (*Block) cmd()  {}

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

	ParenExpr struct {
		X Expr
	}

	UnaryExpr struct {
		Op lexer.Tok
		X  Expr
	}
)

func (*BinaryExpr) node() {}
func (*ParenExpr) node()  {}
func (*UnaryExpr) node()  {}

func (*BinaryExpr) expr() {}
func (*ParenExpr) expr()  {}
func (*UnaryExpr) expr()  {}

type (
	Lit interface {
		lit()
		Expr
	}

	Bool struct {
		Val bool
	}

	F64 struct {
		Val float64
	}

	I64 struct {
		Val int64
	}

	String struct {
		Val string
	}
)

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
