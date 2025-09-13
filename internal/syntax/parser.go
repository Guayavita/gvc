package syntax

import (
	"fmt"

	"jmpeax.com/guayavita/gvc/internal/diag"
)

type ParseError struct {
	Msg string
	Pos Pos
}

func (e *ParseError) Error() string {
	if e.Pos.Line > 0 {
		return fmt.Sprintf("parse error at %d:%d: %s", e.Pos.Line, e.Pos.Col, e.Msg)
	}
	return "parse error: " + e.Msg
}

// Diagnostic converts a ParseError into a diag.Diagnostic for pretty printing.
func (e *ParseError) Diagnostic(file string) diag.Diagnostic {
	start := diag.Position{File: file, Line: e.Pos.Line, Column: e.Pos.Col}
	return diag.Diagnostic{
		Severity: diag.Error,
		Message:  e.Msg,
		Span:     diag.Span{Start: start, End: start},
	}
}

type Parser struct {
	lex   *Lexer
	peek  *Token
	peekE error
}

func NewParser(lx *Lexer) *Parser {
	return &Parser{lex: lx}
}

func (p *Parser) next() (Token, error) {
	if p.peek != nil {
		tok := *p.peek
		err := p.peekE
		p.peek = nil
		p.peekE = nil
		return tok, err
	}
	return p.lex.NextToken()
}

func (p *Parser) expect(tt TokenType) (Token, error) {
	tok, err := p.next()
	if err != nil {
		return tok, err
	}
	if tok.Type != tt {
		return tok, &ParseError{
			Msg: fmt.Sprintf("expected %s, got %s (%q)", tt, tok.Type, tok.Lit),
			Pos: tok.Pos,
		}
	}
	return tok, nil
}

// ParsePackage parses the input according to the grammar:
//
//	PackageDecl -> "package" IDENT EOF
func (p *Parser) ParsePackage() (*PackageDecl, error) {
	tPkg, err := p.expect(PACKAGE)
	if err != nil {
		return nil, err
	}
	tName, err := p.expect(IDENT)
	if err != nil {
		return nil, err
	}

	// Require EOF (lexer tolerates trailing spaces/comments)
	t, err := p.next()
	if err != nil {
		return nil, err
	}
	if t.Type != EOF {
		return nil, &ParseError{
			Msg: fmt.Sprintf("unexpected token after package declaration: %s (%q)", t.Type, t.Lit),
			Pos: t.Pos,
		}
	}

	return &PackageDecl{
		Name: tName.Lit,
		Pos:  tPkg.Pos,
	}, nil
}
