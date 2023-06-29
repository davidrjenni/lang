// Copyright (c) 2023 David Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parser // import "davidrjenni.io/lang/parser"

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"

	"davidrjenni.io/lang/ast"
	"davidrjenni.io/lang/internal/errors"
	"davidrjenni.io/lang/lexer"
)

func ParseFile(filename string) (*ast.Block, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Parse(f, filename)
}

func Parse(r io.Reader, filename string) (*ast.Block, error) {
	l, err := lexer.New(r, filename)
	if err != nil {
		return nil, err
	}
	p := &parser{l: l}
	p.next()
	n := p.parseBlock()
	return n, p.errs.Err()
}

type parser struct {
	l    *lexer.Lexer
	errs errors.Errors

	// lookahead
	pos lexer.Pos
	lit string
	tok lexer.Tok
}

// Block -> "{" Cmd { Cmd } "}" .
func (p *parser) parseBlock() *ast.Block {
	var b ast.Block
	b.StartPos = p.expect(lexer.LeftBrace)
	for p.in(lexer.Assert) {
		b.Cmds = append(b.Cmds, p.parseCmd())
	}
	b.EndPos = p.expect(lexer.RightBrace)
	return &b
}

// Cmd -> "assert" Expr ";" .
func (p *parser) parseCmd() ast.Cmd {
	pos := p.expect(lexer.Assert)
	x := p.parseExpr()
	end := p.expect(lexer.Semicolon)
	return &ast.Assert{X: x, StartPos: pos, EndPos: end}
}

// Expr -> UnaryExpr { BinOp UnaryExpr } .
func (p *parser) parseExpr() ast.Expr {
	var (
		parseMul = func() ast.Expr { return p.parseBinaryExpr(p.parseUnaryExpr, lexer.Multiply, lexer.Divide) }
		parseAdd = func() ast.Expr { return p.parseBinaryExpr(parseMul, lexer.Plus, lexer.Minus) }
		parseRel = func() ast.Expr { return p.parseBinaryExpr(parseAdd, relOps...) }
		parseAnd = func() ast.Expr { return p.parseBinaryExpr(parseRel, lexer.And) }
		parseOr  = func() ast.Expr { return p.parseBinaryExpr(parseAnd, lexer.Or) }
	)
	return p.parseBinaryExpr(parseOr, lexer.Implies)
}

func (p *parser) parseBinaryExpr(parse func() ast.Expr, ops ...lexer.Tok) ast.Expr {
	x := parse()
	for p.in(ops...) {
		op := p.tok
		p.next()
		y := parse()
		x = &ast.BinaryExpr{LHS: x, Op: op, RHS: y}
	}
	return x
}

// UnaryExpr -> [ "-" | "~" ] UnaryExprExpr | PrimaryExpr .
func (p *parser) parseUnaryExpr() ast.Expr {
	switch p.tok {
	case lexer.Minus, lexer.Not:
		op := p.tok
		pos := p.expect(lexer.Minus, lexer.Not)
		x := p.parseUnaryExpr()
		return &ast.UnaryExpr{X: x, Op: op, StartPos: pos}
	default:
		return p.parsePrimaryExpr()
	}
}

// PrimaryExpr -> "(" Expr ")" | F64Lit | I64Lit | StringLit | "true" | "false" .
func (p *parser) parsePrimaryExpr() ast.Expr {
	switch p.tok {
	case lexer.F64Lit:
		lit := p.lit
		pos := p.expect(lexer.F64Lit)
		end := pos.Shift(len(lit))
		val, err := strconv.ParseFloat(lit, 64)
		if err != nil {
			panic(fmt.Sprintf("cannot convert f64: %v", err))
		}
		return &ast.F64{Val: val, StartPos: pos, EndPos: end}
	case lexer.I64Lit:
		lit := p.lit
		pos := p.expect(lexer.I64Lit)
		end := pos.Shift(len(lit))
		val, err := strconv.ParseInt(lit, 10, 0)
		if err != nil {
			panic(fmt.Sprintf("cannot convert i64: %v", err))
		}
		return &ast.I64{Val: val, StartPos: pos, EndPos: end}
	case lexer.LeftParen:
		pos := p.expect(lexer.LeftParen)
		x := p.parseExpr()
		end := p.expect(lexer.RightParen)
		return &ast.ParenExpr{X: x, StartPos: pos, EndPos: end}
	case lexer.StringLit:
		lit := p.lit
		pos := p.expect(lexer.StringLit)
		end := pos.Shift(len(lit))
		return &ast.String{Val: lit, StartPos: pos, EndPos: end}
	case lexer.True:
		lit := p.lit
		pos := p.expect(lexer.True)
		end := pos.Shift(len(lit))
		return &ast.Bool{Val: true, StartPos: pos, EndPos: end}
	case lexer.False:
		lit := p.lit
		pos := p.expect(lexer.False)
		end := pos.Shift(len(lit))
		return &ast.Bool{Val: false, StartPos: pos, EndPos: end}
	default:
		p.errs.Append("unexpected %s", p.lit)
		return nil
	}
}

func (p *parser) next() {
	for {
		var err error
		p.pos, p.tok, p.lit, err = p.l.Read()
		if err != nil {
			p.errs.Append("syntax error: %v", err)
			continue
		}
		if p.tok == lexer.Illegal {
			p.errs.Append("syntax error: %s", p.lit)
			continue
		}
		break
	}
}

func (p *parser) got(toks ...lexer.Tok) bool {
	if p.in(toks...) {
		p.next()
		return true
	}
	return false
}

func (p *parser) in(toks ...lexer.Tok) bool {
	for _, tok := range toks {
		if p.tok == tok {
			return true
		}
	}
	return false
}

func (p *parser) expect(toks ...lexer.Tok) lexer.Pos {
	pos := p.pos
	if p.got(toks...) {
		return pos
	}
	if len(toks) == 1 {
		p.errs.Append("unexpected %s, expected %s", p.lit, toks[0])
		return pos
	}
	var b bytes.Buffer
	for i, tok := range toks {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(tok.String())
	}
	p.errs.Append("unexpected %s, expected one of %s", p.lit, b.String())
	return pos
}

var relOps = []lexer.Tok{
	lexer.Less,
	lexer.LessEq,
	lexer.Equal,
	lexer.NotEqual,
	lexer.GreaterEq,
	lexer.Greater,
	lexer.In,
	lexer.Is,
}
