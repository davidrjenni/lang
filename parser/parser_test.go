// Copyright (c) 2023 David Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parser_test

import (
	"bytes"
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"davidrjenni.io/lang/ast"
	"davidrjenni.io/lang/parser"
)

var update = flag.Bool("update", false, "update golden files")

func TestParse(t *testing.T) {
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

	var actual bytes.Buffer
	ast.Dump(&actual, n)

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
