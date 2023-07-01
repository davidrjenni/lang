// Copyright (c) 2023 David Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ir_test

import (
	"bytes"
	"flag"
	"io/ioutil"
	"path/filepath"
	"testing"

	"davidrjenni.io/lang/ir"
	"davidrjenni.io/lang/parser"
	"davidrjenni.io/lang/types"
)

var update = flag.Bool("update", false, "update golden files")

func TestTranslate(t *testing.T) {
	filename := filepath.Join("test-fixtures", "input.l")
	b, _, err := parser.ParseFile(filename)
	if err != nil {
		t.Fatalf("cannot parse file: %v", err)
	}

	info, err := types.Check(b)
	if err != nil {
		t.Fatalf("%v", err)
	}

	id := ir.Pass(func(s ir.Seq) ir.Seq { return s })

	passes := [...]struct {
		filename string
		pass     ir.Pass
	}{
		{filename: "input.golden", pass: id},
		{filename: "input.loads.golden", pass: ir.Loads},
	}

	for _, p := range passes {
		seq := ir.Translate(b, info, p.pass)
		cmpGolden(t, seq, p.filename, *update)
	}
}

func cmpGolden(t *testing.T, seq ir.Seq, filename string, update bool) {
	var actual bytes.Buffer
	ir.Dump(&actual, seq)

	golden := filepath.Join("test-fixtures", filename)
	if update {
		if err := ioutil.WriteFile(golden, actual.Bytes(), 0644); err != nil {
			t.Fatalf("cannot update golden file: %v", err)
		}
	}

	expected, err := ioutil.ReadFile(golden)
	if err != nil {
		t.Fatalf("cannot read golden file: %v", err)
	}

	if !bytes.Equal(actual.Bytes(), expected) {
		t.Errorf("%s: expected\n%s\ngot\n%s\n", filename, string(expected), actual.String())
	}
}
