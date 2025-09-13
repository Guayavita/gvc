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
