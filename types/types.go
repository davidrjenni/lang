// Copyright (c) 2023 David Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types // import "davidrjenni.io/lang/types"

import "fmt"

type Type interface {
	String() string
	typ()
}

func Equal(t, u Type) bool {
	switch t := t.(type) {
	case *Bool:
		_, ok := u.(*Bool)
		return ok
	case *F64:
		_, ok := u.(*F64)
		return ok
	case *I64:
		_, ok := u.(*I64)
		return ok
	case *String:
		_, ok := u.(*String)
		return ok
	default:
		panic(fmt.Sprintf("unexpected type %T", t))
	}
}

type (
	Bool   struct{}
	F64    struct{}
	I64    struct{}
	String struct{}
)

func (*Bool) String() string   { return "bool" }
func (*F64) String() string    { return "f64" }
func (*I64) String() string    { return "i64" }
func (*String) String() string { return "string" }

func (*Bool) typ()   {}
func (*F64) typ()    {}
func (*I64) typ()    {}
func (*String) typ() {}
