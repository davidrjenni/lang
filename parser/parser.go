// Copyright (c) 2023 David Jenni. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parser // import "davidrjenni.io/lang/parser"

import (
	"bytes"
	"io"
	"os"

	"davidrjenni.io/lang/ast"
	"davidrjenni.io/lang/internal/errors"
	"davidrjenni.io/lang/lexer"
)

func ParseFile(filename string) (*ast.Block, []*ast.Comment, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()
	return Parse(f, filename)
}

func Parse(r io.Reader, filename string) (*ast.Block, []*ast.Comment, error) {
	l, err := lexer.New(r, filename)
	if err != nil {
		return nil, nil, err
	}
	p := &parser{l: l}
	p.next()
	n := p.parseBlock()
	return n, p.comments, p.errs.Err()
}

type parser struct {
	l        *lexer.Lexer
	errs     errors.Errors
	comments []*ast.Comment

	// lookahead
	pos lexer.Pos
	lit string
	tok lexer.Tok
}

// ------- Declarations -------

// VarDecl -> "let" Ident ":=" Expr ";" .
func (p *parser) parseVarDecl() *ast.VarDecl {
	pos := p.expect(lexer.Let)
	ident := p.parseIdent()
	p.expect(lexer.Define)
	x := p.parseExpr()
	end := p.expect(lexer.Semicolon)
	return &ast.VarDecl{Ident: ident, X: x, StartPos: pos, EndPos: end}
}

// ------- Types -------

// Types -> Type { "," Type } .
func (p *parser) parseTypes() (ts []ast.Type) {
	ts = append(ts, p.parseType())
	for p.got(lexer.Comma) {
		ts = append(ts, p.parseType())
	}
	return ts
}

// Type -> "bool" | "f64" | Func | "i64" | "string" .
func (p *parser) parseType() ast.Type {
	switch p.tok {
	case lexer.Bool:
		return p.parseScalar(lexer.Bool)
	case lexer.F64:
		return p.parseScalar(lexer.F64)
	case lexer.Func:
		return p.parseFunc()
	case lexer.I64:
		return p.parseScalar(lexer.I64)
	case lexer.String:
		return p.parseScalar(lexer.String)
	default:
		p.errs.Append(p.pos, "unexpected %s", p.lit)
		return nil
	}
}

// Func -> "func" "(" [ Types ] ")" Type .
func (p *parser) parseFunc() *ast.Func {
	pos := p.expect(lexer.Func)
	p.expect(lexer.LeftParen)
	var params []ast.Type
	if p.tok != lexer.RightParen {
		params = p.parseTypes()
	}
	p.expect(lexer.RightParen)
	result := p.parseType()
	return &ast.Func{Params: params, Result: result, StartPos: pos}
}

func (p *parser) parseScalar(tok lexer.Tok) *ast.Scalar {
	lit := p.lit
	pos := p.expect(tok)
	end := pos.Shift(len(lit))
	return &ast.Scalar{Name: lit, StartPos: pos, EndPos: end}
}

// ------- Commands -------

// Block -> "{" Cmd { Cmd } "}" .
func (p *parser) parseBlock() *ast.Block {
	var b ast.Block
	b.StartPos = p.expect(lexer.LeftBrace)
	for p.in(cmdToks...) {
		b.Cmds = append(b.Cmds, p.parseCmd())
	}
	b.EndPos = p.expect(lexer.RightBrace)
	return &b
}

// Cmd -> Assert | Break | Continue | For | If | VarDecl | Return | Assign .
func (p *parser) parseCmd() ast.Cmd {
	switch p.tok {
	case lexer.Assert:
		return p.parseAssert()
	case lexer.Break:
		return p.parseBreak()
	case lexer.Continue:
		return p.parseContinue()
	case lexer.For:
		return p.parseFor()
	case lexer.If:
		return p.parseIf()
	case lexer.Let:
		return p.parseVarDecl()
	case lexer.Return:
		return p.parseReturn()
	case lexer.Set:
		return p.parseAssign()
	default:
		p.errs.Append(p.pos, "unexpected %s", p.lit)
		return nil
	}
}

// Assert -> "assert" Expr ";" .
func (p *parser) parseAssert() *ast.Assert {
	pos := p.expect(lexer.Assert)
	x := p.parseExpr()
	end := p.expect(lexer.Semicolon)
	return &ast.Assert{X: x, StartPos: pos, EndPos: end}
}

// Break -> "break" ";" .
func (p *parser) parseBreak() *ast.Break {
	pos := p.expect(lexer.Break)
	end := p.expect(lexer.Semicolon)
	return &ast.Break{StartPos: pos, EndPos: end}
}

// Continue -> "continue" ";" .
func (p *parser) parseContinue() *ast.Continue {
	pos := p.expect(lexer.Continue)
	end := p.expect(lexer.Semicolon)
	return &ast.Continue{StartPos: pos, EndPos: end}
}

// For -> "for" Expr Block .
func (p *parser) parseFor() *ast.For {
	pos := p.expect(lexer.For)
	x := p.parseExpr()
	b := p.parseBlock()
	return &ast.For{X: x, Block: b, StartPos: pos}
}

// If  -> "if" Expr Block [ "else" ( If | Block ) ] .
func (p *parser) parseIf() *ast.If {
	pos := p.expect(lexer.If)
	x := p.parseExpr()
	b := p.parseBlock()

	var e *ast.Else
	if p.tok == lexer.Else {
		pos := p.expect(lexer.Else)
		var cmd ast.Cmd
		switch p.tok {
		case lexer.LeftBrace:
			cmd = p.parseBlock()
		case lexer.If:
			cmd = p.parseCmd()
		default:
			p.errs.Append(p.pos, "unexpected %s, expected %s or %s", p.lit, lexer.LeftBrace, lexer.If)
		}
		if cmd != nil {
			e = &ast.Else{Cmd: cmd, StartPos: pos}
		}
	}
	return &ast.If{X: x, Block: b, Else: e, StartPos: pos}
}

// Return -> "return" Expr ";" .
func (p *parser) parseReturn() *ast.Return {
	pos := p.expect(lexer.Return)
	x := p.parseExpr()
	end := p.expect(lexer.Semicolon)
	return &ast.Return{X: x, StartPos: pos, EndPos: end}
}

// Assign -> "set" Ident "<-" Expr ";" .
func (p *parser) parseAssign() *ast.Assign {
	pos := p.expect(lexer.Set)
	ident := p.parseIdent()
	p.expect(lexer.Assign)
	x := p.parseExpr()
	end := p.expect(lexer.Semicolon)
	return &ast.Assign{Ident: ident, X: x, StartPos: pos, EndPos: end}
}

// ------- Expressions -------

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

// UnaryExpr -> [ "-" | "~" ] UnaryExpr | PrimaryExpr .
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

// PrimaryExpr -> ParenExpr | F64Lit | FuncLit | I64Lit | Identifier | StringLit | True | False .
func (p *parser) parsePrimaryExpr() ast.Expr {
	switch p.tok {
	case lexer.LeftParen:
		return p.parseParenExpr()
	case lexer.F64Lit:
		return p.parseF64Lit()
	case lexer.Func:
		return p.parseFuncLit()
	case lexer.I64Lit:
		return p.parseI64Lit()
	case lexer.Identifier:
		return p.parseIdent()
	case lexer.StringLit:
		return p.parseStringLit()
	case lexer.True:
		return p.parseTrue()
	case lexer.False:
		return p.parseFalse()
	default:
		p.errs.Append(p.pos, "unexpected %s", p.lit)
		return nil
	}
}

// ParenExpr -> "(" Expr ")" .
func (p *parser) parseParenExpr() *ast.ParenExpr {
	pos := p.expect(lexer.LeftParen)
	x := p.parseExpr()
	end := p.expect(lexer.RightParen)
	return &ast.ParenExpr{X: x, StartPos: pos, EndPos: end}
}

func (p *parser) parseF64Lit() *ast.F64 {
	lit := p.lit
	pos := p.expect(lexer.F64Lit)
	end := pos.Shift(len(lit))
	return &ast.F64{Val: lit, StartPos: pos, EndPos: end}
}

// FuncLit -> "func" "(" [ Fields ] ")" Type Block .
func (p *parser) parseFuncLit() *ast.FuncLit {
	pos := p.expect(lexer.Func)
	p.expect(lexer.LeftParen)
	var params []*ast.Field
	if p.tok != lexer.RightParen {
		params = p.parseFields()
	}
	p.expect(lexer.RightParen)
	result := p.parseType()
	b := p.parseBlock()
	return &ast.FuncLit{Params: params, Result: result, Block: b, StartPos: pos}
}

func (p *parser) parseI64Lit() *ast.I64 {
	lit := p.lit
	pos := p.expect(lexer.I64Lit)
	end := pos.Shift(len(lit))
	return &ast.I64{Val: lit, StartPos: pos, EndPos: end}
}

func (p *parser) parseIdent() *ast.Ident {
	name := p.lit
	pos := p.expect(lexer.Identifier)
	return &ast.Ident{Name: name, StartPos: pos}
}

func (p *parser) parseStringLit() *ast.String {
	lit := p.lit
	pos := p.expect(lexer.StringLit)
	end := pos.Shift(len(lit))
	return &ast.String{Val: lit, StartPos: pos, EndPos: end}
}

// True -> "true" .
func (p *parser) parseTrue() *ast.Bool {
	lit := p.lit
	pos := p.expect(lexer.True)
	end := pos.Shift(len(lit))
	return &ast.Bool{Val: lit, StartPos: pos, EndPos: end}
}

// False -> "false" .
func (p *parser) parseFalse() *ast.Bool {
	lit := p.lit
	pos := p.expect(lexer.False)
	end := pos.Shift(len(lit))
	return &ast.Bool{Val: lit, StartPos: pos, EndPos: end}
}

// ------- Others -------

// Fields -> Field { "," Field } .
// Field  -> Ident Type .
func (p *parser) parseFields() (fs []*ast.Field) {
	id := p.parseIdent()
	typ := p.parseType()
	fs = append(fs, &ast.Field{Ident: id, Type: typ})
	for p.got(lexer.Comma) {
		id = p.parseIdent()
		typ = p.parseType()
		fs = append(fs, &ast.Field{Ident: id, Type: typ})
	}
	return fs
}

// ------- Helpers -------

func (p *parser) next() {
	for {
		var err error
		p.pos, p.tok, p.lit, err = p.l.Read()
		if err != nil {
			p.errs.Append(p.pos, "syntax error: %v", err)
			continue
		}
		if p.tok == lexer.Illegal {
			p.errs.Append(p.pos, "syntax error: %s", p.lit)
			continue
		}
		if p.tok == lexer.Comment {
			p.comments = append(p.comments, &ast.Comment{
				Text:     p.lit,
				StartPos: p.pos,
				EndPos:   p.pos.Shift(len(p.lit)),
			})
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
		p.errs.Append(p.pos, "unexpected %s, expected %s", p.lit, toks[0])
		return pos
	}
	var b bytes.Buffer
	for i, tok := range toks {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(tok.String())
	}
	p.errs.Append(p.pos, "unexpected %s, expected one of %s", p.lit, b.String())
	return pos
}

var (
	relOps = []lexer.Tok{
		lexer.Less,
		lexer.LessEq,
		lexer.Equal,
		lexer.NotEqual,
		lexer.GreaterEq,
		lexer.Greater,
		lexer.In,
		lexer.Is,
	}

	cmdToks = []lexer.Tok{
		lexer.Assert,
		lexer.Break,
		lexer.Continue,
		lexer.For,
		lexer.If,
		lexer.Let,
		lexer.Return,
		lexer.Set,
	}
)
