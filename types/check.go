// Copyright (c) 2023 David Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types // import "davidrjenni.io/lang/types"

import (
	"fmt"

	"davidrjenni.io/lang/ast"
	"davidrjenni.io/lang/internal/errors"
	"davidrjenni.io/lang/lexer"
)

func Check(n ast.Node) error {
	c := &checker{}
	c.check(n)
	return c.errs.Err()
}

type checker struct {
	errs errors.Errors
}

func (c *checker) check(n ast.Node) {
	switch n := n.(type) {
	case *ast.Assert:
		t, ok := c.checkExpr(n.X)
		if !ok {
			return
		}
		if _, ok := t.(*Bool); !ok {
			c.errorf("expr must be of type bool, got %s", t)
		}
	case *ast.Block:
		for _, cmd := range n.Cmds {
			c.check(cmd)
		}
	default:
		panic(fmt.Sprintf("unexpected type %T", n))
	}
}

func (c *checker) checkExpr(x ast.Expr) (Type, bool) {
	switch x := x.(type) {
	case *ast.BinaryExpr:
		return c.checkBinaryExpr(x)
	case *ast.Bool:
		return &Bool{}, true
	case *ast.I64:
		return &I64{}, true
	case *ast.F64:
		return &F64{}, true
	case *ast.ParenExpr:
		return c.checkExpr(x.X)
	case *ast.String:
		return &String{}, true
	case *ast.UnaryExpr:
		return c.checkUnaryExpr(x)
	default:
		panic(fmt.Sprintf("unexpected type %T", x))
	}
}

func (c *checker) checkBinaryExpr(x *ast.BinaryExpr) (Type, bool) {
	lhs, ok := c.checkExpr(x.LHS)
	if !ok {
		return nil, false
	}
	rhs, ok := c.checkExpr(x.RHS)
	if !ok {
		return nil, false
	}
	if !Equal(lhs, rhs) {
		c.errorf("cannot apply %s to operands of types %s and %s", x.Op, lhs, rhs)
		return nil, false
	}

	switch t := lhs.(type) {
	case *Bool:
		if lexer.And <= x.Op && x.Op <= lexer.Implies {
			return &Bool{}, true
		}
		if lexer.Equal == x.Op || lexer.NotEqual == x.Op {
			return &Bool{}, true
		}
	case *I64:
		if lexer.Plus <= x.Op && x.Op <= lexer.Divide {
			return &I64{}, true
		}
		if lexer.Less <= x.Op && x.Op <= lexer.GreaterEq {
			return &Bool{}, true
		}
	case *F64:
		if lexer.Plus <= x.Op && x.Op <= lexer.Divide {
			return &F64{}, true
		}
		if lexer.Less <= x.Op && x.Op <= lexer.GreaterEq {
			return &Bool{}, true
		}
	case *String:
		if x.Op == lexer.Plus {
			return &String{}, true
		}
		if lexer.Less <= x.Op && x.Op <= lexer.GreaterEq {
			return &Bool{}, true
		}
	default:
		panic(fmt.Sprintf("unexpected type %T", t))
	}

	c.errorf("cannot apply %s to operands of types %s and %s", x.Op, lhs, rhs)
	return nil, false
}

func (c *checker) checkUnaryExpr(x *ast.UnaryExpr) (Type, bool) {
	t, ok := c.checkExpr(x.X)
	if !ok {
		return nil, false
	}

	switch t := t.(type) {
	case *Bool:
		if x.Op != lexer.Not {
			c.errorf("cannot apply %s to expr of type %s", x.Op, t)
			return nil, false
		}
		return &Bool{}, true
	case *I64:
		if x.Op != lexer.Minus {
			c.errorf("cannot apply %s to expr of type %s", x.Op, t)
			return nil, false
		}
		return &I64{}, true
	case *F64:
		if x.Op != lexer.Minus {
			c.errorf("cannot apply %s to expr of type %s", x.Op, t)
			return nil, false
		}
		return &F64{}, true
	default:
		c.errorf("cannot apply %s to expr of type %s", x.Op, t)
		return nil, false
	}
}

func (c *checker) errorf(format string, args ...interface{}) {
	c.errs.Append(format, args...)
}
