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

	True  // true
	False // false
)

// keywords map all keywords to their corresponding token.
var keywords = map[string]Tok{
	"true":  True,
	"false": False,
}
