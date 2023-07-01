// Copyright (c) 2023 David Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compiler_test

import (
	"bytes"
	"flag"
	"io/ioutil"
	"path/filepath"
	"testing"

	"davidrjenni.io/lang/compiler"
	"davidrjenni.io/lang/ir"
	"davidrjenni.io/lang/parser"
	"davidrjenni.io/lang/types"
)

var update = flag.Bool("update", false, "update golden files")

func TestCompile(t *testing.T) {
	filename := filepath.Join("test-fixtures", "input.l")
	b, err := parser.ParseFile(filename)
	if err != nil {
		t.Fatalf("cannot parse file: %v", err)
	}

	info, err := types.Check(b)
	if err != nil {
		t.Fatalf("%v", err)
	}

	n := ir.Translate(b, info)

	var actual bytes.Buffer
	compiler.Compile(&actual, filename, n)

	golden := filepath.Join("test-fixtures", "input.golden")
	if *update {
		if err := ioutil.WriteFile(golden, actual.Bytes(), 0644); err != nil {
			t.Fatalf("cannot update golden file: %v", err)
		}
	}

	expected, err := ioutil.ReadFile(golden)
	if err != nil {
		t.Fatalf("cannot read golden file: %v", err)
	}

	if !bytes.Equal(actual.Bytes(), expected) {
		t.Fatalf("expected\n%s\ngot\n%s\n", string(expected), actual.String())
	}
}
