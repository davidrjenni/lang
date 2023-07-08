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
	Uses  map[*ast.Ident]*Object
	Types map[ast.Expr]*Object
}

type scope struct {
	parent *scope

	objects map[string]*Object
	func_   *Func
	inFor   bool
}

func (s *scope) enter() *scope {
	return &scope{
		parent:  s,
		objects: make(map[string]*Object),
		func_:   s.func_,
		inFor:   s.inFor,
	}
}

func (s *scope) lookup(name string) (*Object, bool) {
	obj, ok := s.objects[name]
	if !ok && s.parent != nil {
		return s.parent.lookup(name)
	}
	return obj, ok
}
