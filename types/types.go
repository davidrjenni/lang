// Copyright (c) 2023 David Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types // import "davidrjenni.io/lang/types"

import (
	"bytes"
	"fmt"
)

type Type interface {
	Size() int
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
	case *Func:
		f, ok := u.(*Func)
		if !ok {
			return false
		}
		if len(t.Params) != len(f.Params) {
			return false
		}
		for i, p := range t.Params {
			if !Equal(p, f.Params[i]) {
				return false
			}
		}
		return Equal(t.Result, f.Result)
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

type Func struct {
	Params []Type
	Result Type
}

func (f *Func) String() string {
	b := bytes.NewBufferString("func(")
	for i, p := range f.Params {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(p.String())
	}
	b.WriteString(") ")
	b.WriteString(f.Result.String())
	return b.String()
}

type (
	Bool   struct{}
	F64    struct{}
	I64    struct{}
	String struct{}
)

func (*Bool) Size() int   { return 1 }
func (*F64) Size() int    { return 8 }
func (*Func) Size() int   { return 8 }
func (*I64) Size() int    { return 8 }
func (*String) Size() int { return 8 }

func (*Bool) String() string   { return "bool" }
func (*F64) String() string    { return "f64" }
func (*I64) String() string    { return "i64" }
func (*String) String() string { return "string" }

func (*Bool) typ()   {}
func (*F64) typ()    {}
func (*Func) typ()   {}
func (*I64) typ()    {}
func (*String) typ() {}
