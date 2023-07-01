// Copyright (c) 2023 David Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types_test

import (
	"path/filepath"
	"testing"

	"davidrjenni.io/lang/parser"
	"davidrjenni.io/lang/types"
)

func TestCheck(t *testing.T) {
	filename := filepath.Join("test-fixtures", "input.l")
	n, err := parser.ParseFile(filename)
	if err != nil {
		t.Fatalf("cannot parse file: %v", err)
	}

	if _, err := types.Check(n); err != nil {
		t.Fatalf("%v", err)
	}
}
