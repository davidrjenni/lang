// Copyright (c) 2023 David Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lexer // import "davidrjenni.io/lang/lexer"

import (
	"bufio"
	"io"
	"unicode"
)

// eof is the sentinel value to represent "end of file" as a rune.
const eof = -1

// Lexer holds the state of the lexical analyzer.
type Lexer struct {
	r  *bufio.Reader // source reader
	ch rune          // current rune
}

func New(r io.Reader, filename string) (*Lexer, error) {
	l := &Lexer{r: bufio.NewReader(r)}

	// Initialize the current rune.
	err := l.next()
	return l, err
}

// Read reads the next token in the reader and returns its
// token type and literal or returns an error. For unknown
// tokens, Illegal is returned as token type. To indicate
// the end of the input, EOF is returned. If an error is
// returned, all other return values are invalid.
func (l *Lexer) Read() (tok Tok, lit string, err error) {
	// Skip the spaces.
	for unicode.IsSpace(l.ch) {
		if err := l.next(); err != nil {
			return tok, lit, err
		}
	}
	lit = string(l.ch)

	ch := l.ch

	// Advance l.ch to the next rune, ch is the previous one.
	if err := l.next(); err != nil {
		return tok, lit, err
	}

	switch ch {
	case eof:
		tok, lit = EOF, "EOF"
	default:
		tok = Illegal
	}

	return tok, lit, nil
}

func (l *Lexer) next() error {
	ch, _, err := l.r.ReadRune()
	switch {
	case err == io.EOF:
		ch = eof
	case err != nil:
		return err
	}

	l.ch = ch
	return nil
}