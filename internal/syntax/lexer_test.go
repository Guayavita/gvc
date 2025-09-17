package syntax

import (
	"testing"
)

func TestLexer_BasicTokens(t *testing.T) {
	input := `package test
// a comment line
/* block
comment */
fun main() : none {
    def x: i32 = 10
    def y = "hi\\n"
    if x >= 10 && y != "" {
        print(x + 2)
    }
}`
	l := NewLexer(input, "<mem>")

	// Iterate until EOF and count some categories
	var (
		idents, ints, strings, ops, keywords int
	)
	for {
		tok := l.NextToken()
		if tok.Kind == EOF {
			break
		}
		// comments should be skipped by lexer
		switch tok.Kind {
		case IDENT:
			idents++
		case INT:
			ints++
		case STRING:
			strings++
		case PLUS, MINUS, MUL, DIV, MOD, ASSIGN, EQ, NE, LT, LE, GT, GE, AND, OR, NOT:
			ops++
		case PACKAGE, FUN, IF, ELSE, RETURN, WHILE, FOR, IN, NONE:
			keywords++
		}
	}
	if idents == 0 {
		t.Fatalf("expected some identifiers, got 0")
	}
	if ints == 0 {
		t.Fatalf("expected some ints, got 0")
	}
	if strings == 0 {
		t.Fatalf("expected some strings, got 0")
	}
	if ops == 0 {
		t.Fatalf("expected some operators, got 0")
	}
	if keywords == 0 {
		t.Fatalf("expected some keywords, got 0")
	}
}
