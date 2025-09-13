package syntax

import "testing"

func TestParseFile_DefWithEmptyCall(t *testing.T) {
	src := "package test\n\n def R = DoIt()\n"
	lx := NewLexerFromString(src)
	ps := NewParser(lx)
	f, err := ps.ParseFile()
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if len(f.Definitions) != 1 {
		t.Fatalf("expected 1 definition, got %d", len(f.Definitions))
	}
	def := f.Definitions[0]
	call, ok := def.Value.(CallExpr)
	if !ok {
		t.Fatalf("expected CallExpr as value")
	}
	if call.Name != "DoIt" {
		t.Fatalf("expected call name DoIt, got %q", call.Name)
	}
	if len(call.Args) != 0 {
		t.Fatalf("expected 0 args, got %d", len(call.Args))
	}
}

func TestParseFile_DefWithCallArgs(t *testing.T) {
	src := "package test\n\n def X = Fn(a,1,b2)\n"
	lx := NewLexerFromString(src)
	ps := NewParser(lx)
	f, err := ps.ParseFile()
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if len(f.Definitions) != 1 {
		t.Fatalf("expected 1 definition, got %d", len(f.Definitions))
	}
	def := f.Definitions[0]
	call, ok := def.Value.(CallExpr)
	if !ok {
		t.Fatalf("expected CallExpr as value")
	}
	if call.Name != "Fn" {
		t.Fatalf("expected call name Fn, got %q", call.Name)
	}
	if len(call.Args) != 3 {
		t.Fatalf("expected 3 args, got %d", len(call.Args))
	}
	if _, ok := call.Args[0].(IdentExpr); !ok {
		t.Fatalf("arg0 should be IdentExpr")
	}
	if _, ok := call.Args[1].(NumberExpr); !ok {
		t.Fatalf("arg1 should be NumberExpr")
	}
	if _, ok := call.Args[2].(IdentExpr); !ok {
		t.Fatalf("arg2 should be IdentExpr")
	}
}
