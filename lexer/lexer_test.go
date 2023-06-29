// Copyright (c) 2023 David Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lexer_test

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"davidrjenni.io/lang/lexer"
)

var update = flag.Bool("update", false, "update golden files")

func TestLexer(t *testing.T) {
	filename := "input.l"
	path := filepath.Join("test-fixtures", filename)
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("cannot open file: %v", err)
	}
	defer f.Close()

	l, err := lexer.New(f, filename)
	if err != nil {
		t.Fatalf("cannot initialize lexer: %v", err)
	}

	var actual bytes.Buffer
	for {
		pos, tok, lit, err := l.Read()
		if err != nil {
			t.Fatalf("cannot read from lexer: %v", err)
		}
		fmt.Fprintf(&actual, "%s: %v | %v\n", pos, tok, lit)
		if tok == lexer.EOF {
			break
		}
	}

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
