// Copyright (c) 2023 David Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ast_test

import "davidrjenni.io/lang/ast"

var (
	_ ast.Node = &ast.Assert{}
	_ ast.Node = &ast.Assign{}
	_ ast.Node = &ast.Block{}
	_ ast.Node = &ast.BinaryExpr{}
	_ ast.Node = &ast.Bool{}
	_ ast.Node = &ast.Break{}
	_ ast.Node = &ast.Comment{}
	_ ast.Node = &ast.Continue{}
	_ ast.Node = &ast.Else{}
	_ ast.Node = &ast.F64{}
	_ ast.Node = &ast.Field{}
	_ ast.Node = &ast.For{}
	_ ast.Node = &ast.Func{}
	_ ast.Node = &ast.FuncLit{}
	_ ast.Node = &ast.I64{}
	_ ast.Node = &ast.Ident{}
	_ ast.Node = &ast.If{}
	_ ast.Node = &ast.ParenExpr{}
	_ ast.Node = &ast.Return{}
	_ ast.Node = &ast.Scalar{}
	_ ast.Node = &ast.String{}
	_ ast.Node = &ast.UnaryExpr{}
	_ ast.Node = &ast.VarDecl{}

	_ ast.Decl = &ast.VarDecl{}

	_ ast.Type = &ast.Func{}
	_ ast.Type = &ast.Scalar{}

	_ ast.Cmd = &ast.Assert{}
	_ ast.Cmd = &ast.Assign{}
	_ ast.Cmd = &ast.Block{}
	_ ast.Cmd = &ast.Break{}
	_ ast.Cmd = &ast.Continue{}
	_ ast.Cmd = &ast.For{}
	_ ast.Cmd = &ast.If{}
	_ ast.Cmd = &ast.Return{}
	_ ast.Cmd = &ast.VarDecl{}

	_ ast.Expr = &ast.BinaryExpr{}
	_ ast.Expr = &ast.ParenExpr{}
	_ ast.Expr = &ast.UnaryExpr{}
	_ ast.Expr = &ast.Bool{}
	_ ast.Expr = &ast.F64{}
	_ ast.Expr = &ast.FuncLit{}
	_ ast.Expr = &ast.I64{}
	_ ast.Expr = &ast.Ident{}
	_ ast.Expr = &ast.String{}

	_ ast.Lit = &ast.Bool{}
	_ ast.Lit = &ast.F64{}
	_ ast.Lit = &ast.FuncLit{}
	_ ast.Lit = &ast.I64{}
	_ ast.Lit = &ast.String{}
)
