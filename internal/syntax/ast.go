package syntax

// Pos represents a position in the source input.
type Pos struct {
	Offset int // byte offset from start of input
	Line   int // 1-based
	Col    int // 1-based, in runes
}

// PackageDecl represents: "package IDENT".
type PackageDecl struct {
	Name string
	Pos  Pos // position of the 'package' keyword
}

// File represents a parsed file with an optional package declaration and zero or more definitions.
type File struct {
	Package     *PackageDecl
	Definitions []Definition
}

type Definition struct {
	Name  string
	Value Expr
	Pos   Pos // position of the 'def' keyword
}

// Expr is the interface implemented by all expression nodes.
type Expr interface {
	isExpr()
}

// IdentExpr is an expression that references another identifier.
type IdentExpr struct {
	Name string
	Pos  Pos
}

func (IdentExpr) isExpr() {}

// NumberExpr is an expression representing a numeric literal (int or float).
type NumberExpr struct {
	Value string // as written in source (underscores preserved)
	Pos   Pos
}

func (NumberExpr) isExpr() {}
