// Copyright (c) 2023 David Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types // import "davidrjenni.io/lang/types"

import "davidrjenni.io/lang/ast"

type Object struct {
	Node ast.Node
	Type Type
}

type Info struct {
	Types map[ast.Expr]*Object
}
