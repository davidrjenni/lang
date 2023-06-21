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

// New creates a new lexer which reads source code from the given reader.
// If the lexer cannot read from the reader, an error is returned.
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

	if unicode.IsLetter(l.ch) {
		return l.scanKeyword()
	}
	if unicode.IsDigit(l.ch) {
		return l.scanNumber()
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

	case '(':
		tok = LeftParen
	case ')':
		tok = RightParen
	case '[':
		tok = LeftBracket
	case ']':
		tok = RightBracket
	case ',':
		tok = Comma
	case ';':
		tok = Semicolon

	case '+':
		tok = Plus
	case '-':
		tok = Minus
	case '*', '·':
		tok = Multiply
	case '/', '÷':
		tok = Divide
	case '&', '∧':
		tok = And
	case '|', '∨':
		tok = Or
	case '⟹':
		tok = Implies
	case '<':
		if tok = Less; l.ch == '=' {
			tok, lit = LessEq, "<="
			if err := l.next(); err != nil {
				return tok, lit, err
			}
		}
	case '≤':
		tok = LessEq
	case '=':
		if tok = Equal; l.ch == '>' {
			tok, lit = Implies, "=>"
			if err := l.next(); err != nil {
				return tok, lit, err
			}
		}
	case '#', '≠':
		tok = NotEqual
	case '>':
		if tok = Greater; l.ch == '=' {
			tok, lit = GreaterEq, ">="
			if err := l.next(); err != nil {
				return tok, lit, err
			}
		}
	case '≥':
		tok = GreaterEq
	case '∈':
		tok = In
	case '~', '¬':
		tok = Not

	case '"':
		return l.scanString()

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

func (l *Lexer) scanKeyword() (Tok, string, error) {
	var buf []rune
	for unicode.IsLetter(l.ch) || unicode.IsDigit(l.ch) {
		buf = append(buf, l.ch)
		if err := l.next(); err != nil {
			return Illegal, "", err
		}
	}

	lit := string(buf)
	if t, ok := keywords[lit]; ok {
		return t, lit, nil
	}
	return Illegal, lit, nil
}

func (l *Lexer) scanNumber() (Tok, string, error) {
	var (
		buf []rune
		tok = I64Lit
	)
	for unicode.IsDigit(l.ch) || l.ch == '_' || l.ch == '.' {
		if l.ch != '_' {
			buf = append(buf, l.ch)
		}
		if l.ch == '.' {
			if tok == F64Lit {
				tok = Illegal
			} else {
				tok = F64Lit
			}
		}
		if err := l.next(); err != nil {
			return tok, "", err
		}
	}
	return tok, string(buf), nil
}

func (l *Lexer) scanString() (Tok, string, error) {
	var (
		buf = []rune{'"'}
		tok = StringLit
	)
	for l.ch != '"' {
		buf = append(buf, l.ch)
		if l.ch == '\\' {
			if err := l.next(); err != nil {
				return tok, "", err
			}
			buf = append(buf, l.ch)
			if !l.escape('"') {
				tok = Illegal
			}
		}
		if l.ch == '\n' || l.ch == eof {
			tok = Illegal
			break
		}
		if err := l.next(); err != nil {
			return tok, "", err
		}
	}
	buf = append(buf, l.ch)
	if err := l.next(); err != nil {
		return tok, "", err
	}
	return tok, string(buf), nil
}

func (l *Lexer) escape(q rune) bool {
	switch l.ch {
	case 'n', 't', '\\', q:
		return true
	default:
		return false
	}
}
