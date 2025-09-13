package syntax

import "testing"

func TestLexer_StringAndBraces(t *testing.T) {
	lx := NewLexerFromString("\"hello\" { } [ ] < > : ; . ? ! + - * / % == != <= >= && ||")
	// string
	if tok, _ := lx.NextToken(); tok.Type != STRING || tok.Lit != "hello" {
		t.Fatalf("want STRING 'hello', got %v %q", tok.Type, tok.Lit)
	}
	// {
	if tok, _ := lx.NextToken(); tok.Type != LBRACE {
		t.Fatalf("want LBRACE, got %v", tok.Type)
	}
	// }
	if tok, _ := lx.NextToken(); tok.Type != RBRACE {
		t.Fatalf("want RBRACE, got %v", tok.Type)
	}
	// [
	if tok, _ := lx.NextToken(); tok.Type != LBRACKET {
		t.Fatalf("want LBRACKET, got %v", tok.Type)
	}
	// ]
	if tok, _ := lx.NextToken(); tok.Type != RBRACKET {
		t.Fatalf("want RBRACKET, got %v", tok.Type)
	}
	// <
	if tok, _ := lx.NextToken(); tok.Type != LANGLE {
		t.Fatalf("want LANGLE, got %v", tok.Type)
	}
	// >
	if tok, _ := lx.NextToken(); tok.Type != RANGLE {
		t.Fatalf("want RANGLE, got %v", tok.Type)
	}
	// : ; . ? ! + - * / %
	types := []TokenType{COLON, SEMICOLON, DOT, QUESTION, BANG, PLUS, MINUS, STAR, SLASH, PERCENT}
	for _, tt := range types {
		if tok, _ := lx.NextToken(); tok.Type != tt {
			t.Fatalf("want %v, got %v", tt, tok.Type)
		}
	}
	// == != <= >= && ||
	pairs := []struct {
		tt  TokenType
		lit string
	}{
		{EQEQ, "=="}, {NEQ, "!="}, {LTE, "<="}, {GTE, ">="}, {ANDAND, "&&"}, {OROR, "||"},
	}
	for _, p := range pairs {
		if tok, _ := lx.NextToken(); tok.Type != p.tt || tok.Lit != p.lit {
			t.Fatalf("want %v %q, got %v %q", p.tt, p.lit, tok.Type, tok.Lit)
		}
	}
}
