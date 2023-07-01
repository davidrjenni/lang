// Copyright (c) 2023 David Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ast_test

import "davidrjenni.io/lang/ast"

var (
	_ ast.Node = &ast.Assert{}
	_ ast.Node = &ast.Block{}
	_ ast.Node = &ast.BinaryExpr{}
	_ ast.Node = &ast.Bool{}
	_ ast.Node = &ast.F64{}
	_ ast.Node = &ast.For{}
	_ ast.Node = &ast.I64{}
	_ ast.Node = &ast.If{}
	_ ast.Node = &ast.ParenExpr{}
	_ ast.Node = &ast.String{}
	_ ast.Node = &ast.UnaryExpr{}

	_ ast.Cmd = &ast.Assert{}
	_ ast.Cmd = &ast.Block{}
	_ ast.Cmd = &ast.For{}
	_ ast.Cmd = &ast.If{}

	_ ast.Expr = &ast.BinaryExpr{}
	_ ast.Expr = &ast.ParenExpr{}
	_ ast.Expr = &ast.UnaryExpr{}
	_ ast.Expr = &ast.Bool{}
	_ ast.Expr = &ast.F64{}
	_ ast.Expr = &ast.I64{}
	_ ast.Expr = &ast.String{}

	_ ast.Lit = &ast.Bool{}
	_ ast.Lit = &ast.F64{}
	_ ast.Lit = &ast.I64{}
	_ ast.Lit = &ast.String{}
)
