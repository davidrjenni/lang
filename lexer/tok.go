// Copyright (c) 2023 David Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lexer // import "davidrjenni.io/lang/lexer"

// Tok represents the set of lexical tokens.
type Tok int

//go:generate stringer -type=Tok -linecomment

const (
	// EOF marks the end of file.
	EOF     Tok = iota // EOF
	Illegal            // illegal

	LeftParen    // (
	RightParen   // )
	LeftBracket  // [
	RightBracket // ]

	Comma // ,

	Plus      // +
	Minus     // -
	Multiply  // ·
	Divide    // ÷
	And       // ∧
	Or        // ∨
	Implies   // ⟹
	Less      // <
	LessEq    // ≤
	Equal     // =
	NotEqual  // ≠
	Greater   // >
	GreaterEq // ≥
	In        // ∈
	Is        // is
	Not       // ¬

	True  // true
	False // false

	Bool   // bool
	I64    // i64
	F64    // f64
	String // string
)

// keywords map all keywords to their corresponding token.
var keywords = map[string]Tok{
	"true":  True,
	"false": False,

	"in": In,
	"is": Is,

	"bool":   Bool,
	"i64":    I64,
	"f64":    F64,
	"string": String,
}
