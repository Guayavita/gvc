package syntax

import "testing"

func TestLexer_SimplePackage(t *testing.T) {
	lx := NewLexerFromString("package mypkg\n")
	tok, err := lx.NextToken()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok.Type != PACKAGE {
		t.Fatalf("expected PACKAGE, got %v", tok.Type)
	}

	tok, err = lx.NextToken()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok.Type != IDENT || tok.Lit != "mypkg" {
		t.Fatalf("expected IDENT 'mypkg', got %v %q", tok.Type, tok.Lit)
	}

	tok, err = lx.NextToken()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok.Type != EOF {
		t.Fatalf("expected EOF, got %v", tok.Type)
	}
}

func TestLexer_SkipsWhitespaceAndComments(t *testing.T) {
	src := "\n // comment here\n\tpackage\tname// trailing comment\n\n"
	lx := NewLexerFromString(src)
	tok, _ := lx.NextToken()
	if tok.Type != PACKAGE {
		t.Fatalf("want PACKAGE, got %v", tok.Type)
	}
	tok, _ = lx.NextToken()
	if tok.Type != IDENT || tok.Lit != "name" {
		t.Fatalf("want IDENT 'name', got %v %q", tok.Type, tok.Lit)
	}
	tok, _ = lx.NextToken()
	if tok.Type != EOF {
		t.Fatalf("want EOF, got %v", tok.Type)
	}
}

func TestLexer_IllegalRune(t *testing.T) {
	lx := NewLexerFromString("$")
	tok, err := lx.NextToken()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok.Type != ILLEGAL {
		t.Fatalf("expected ILLEGAL, got %v", tok.Type)
	}
	if tok.Lit != "$" {
		t.Fatalf("expected '$' literal, got %q", tok.Lit)
	}
}
