// Copyright (c) 2023 David Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ir_test

import "davidrjenni.io/lang/ir"

var (
	_ ir.Node = &ir.BinaryExpr{}
	_ ir.Node = ir.Bool(false)
	_ ir.Node = ir.Call("")
	_ ir.Node = ir.CJump("")
	_ ir.Node = ir.F64(0)
	_ ir.Node = ir.I64(0)
	_ ir.Node = ir.Jump("")
	_ ir.Node = ir.Label("")
	_ ir.Node = &ir.Load{}
	_ ir.Node = &ir.Mem{}
	_ ir.Node = &ir.Reg{}
	_ ir.Node = ir.Seq{}
	_ ir.Node = &ir.SeqExpr{}
	_ ir.Node = &ir.Store{}
	_ ir.Node = &ir.UnaryExpr{}

	_ ir.RVal = ir.Bool(false)
	_ ir.RVal = ir.F64(0)
	_ ir.RVal = ir.I64(0)
	_ ir.RVal = &ir.Mem{}
	_ ir.RVal = &ir.Reg{}
	_ ir.RVal = &ir.SeqExpr{}

	_ ir.LVal = &ir.Mem{}
	_ ir.LVal = &ir.Reg{}
)
