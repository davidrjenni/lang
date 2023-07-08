// Copyright (c) 2023 David Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types_test

import (
	"testing"

	"davidrjenni.io/lang/types"
)

func TestEqual(t *testing.T) {
	tests := [...]struct {
		t, u     types.Type
		expected bool
	}{
		{t: &types.Bool{}, u: &types.Bool{}, expected: true},
		{t: &types.I64{}, u: &types.I64{}, expected: true},
		{t: &types.F64{}, u: &types.F64{}, expected: true},
		{t: &types.String{}, u: &types.String{}, expected: true},

		{t: &types.Bool{}, u: &types.I64{}, expected: false},
		{t: &types.Bool{}, u: &types.F64{}, expected: false},
		{t: &types.Bool{}, u: &types.String{}, expected: false},

		{t: &types.I64{}, u: &types.Bool{}, expected: false},
		{t: &types.I64{}, u: &types.F64{}, expected: false},
		{t: &types.I64{}, u: &types.String{}, expected: false},

		{t: &types.F64{}, u: &types.Bool{}, expected: false},
		{t: &types.F64{}, u: &types.I64{}, expected: false},
		{t: &types.F64{}, u: &types.String{}, expected: false},

		{t: &types.String{}, u: &types.Bool{}, expected: false},
		{t: &types.String{}, u: &types.I64{}, expected: false},
		{t: &types.String{}, u: &types.F64{}, expected: false},

		{t: &types.Func{Result: &types.Bool{}}, u: &types.Func{Result: &types.Bool{}}, expected: true},
		{t: &types.Func{Result: &types.Bool{}}, u: &types.Func{Result: &types.String{}}, expected: false},
		{t: &types.Func{Result: &types.Bool{}}, u: &types.Func{Params: []types.Type{&types.String{}}, Result: &types.Bool{}}, expected: false},
		{t: &types.Func{Params: []types.Type{&types.I64{}}, Result: &types.Bool{}}, u: &types.Func{Params: []types.Type{&types.String{}}, Result: &types.Bool{}}, expected: false},
		{t: &types.Func{Result: &types.Bool{}}, u: &types.Bool{}, expected: false},
		{t: &types.Bool{}, u: &types.Func{Result: &types.Bool{}}, expected: false},
	}

	for i, test := range tests {
		actual := types.Equal(test.t, test.u)
		if actual != test.expected {
			t.Errorf("%d: expected %s to be equal to %s", i, test.t, test.u)
		}
	}
}
