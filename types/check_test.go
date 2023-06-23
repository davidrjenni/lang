// Copyright (c) 2023 David Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types_test

import (
	"os"
	"path/filepath"
	"testing"

	"davidrjenni.io/lang/parser"
	"davidrjenni.io/lang/types"
)

func TestCheck(t *testing.T) {
	filename := filepath.Join("test-fixtures", "input.l")
	f, err := os.Open(filename)
	if err != nil {
		t.Fatalf("cannot parse file: %v", err)
	}
	defer f.Close()

	n, err := parser.ParseFile(filename)
	if err != nil {
		t.Fatalf("cannot parse file: %v", err)
	}

	if err = types.Check(n); err != nil {
		t.Fatalf("%v", err)
	}
}
