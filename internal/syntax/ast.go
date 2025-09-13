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

// CallExpr represents a function call expression: NAME '(' args? ')'
// Arguments are limited to identifiers or numeric literals for now.
type CallExpr struct {
	Name string
	Args []Expr
	Pos  Pos // position of the function name
}

func (CallExpr) isExpr() {}

// Additional AST nodes to support extended EBNF grammar

// Program is the root of a compilation unit according to the extended grammar.
type Program struct {
	Package *PackageDecl
	Imports []ImportDecl
	Decls   []TopLevelDecl
}

type ImportDecl struct {
	Path  []string // package.path.segments
	Items []string // optional selective imports inside { }
	Pos   Pos
}

type TopLevelDecl interface{ isTopLevelDecl() }

// ConstDecl corresponds to: [export] def name = expression
type ConstDecl struct {
	Exported bool
	Name     string
	Value    Expr
	Pos      Pos
}

func (ConstDecl) isTopLevelDecl() {}

// TypeDecl corresponds to: [export] type Name = (struct | enum)
type TypeDecl struct {
	Exported bool
	Name     string
	Type     TypeExpr // struct/enum/alias
	Pos      Pos
}

func (TypeDecl) isTopLevelDecl() {}

type StructDecl struct {
	Fields []FieldDecl
	Pos    Pos
}
type FieldDecl struct {
	Name string
	Type TypeExpr
	Pos  Pos
}

type EnumDecl struct {
	Variants []VariantDecl
	Pos      Pos
}
type VariantDecl struct {
	Name  string
	Types []TypeExpr
	Pos   Pos
}

// ImplBlock corresponds to: impl [<T,...>] Type { fun ... }
type ImplBlock struct {
	TypeParams []string
	TypeName   string
	Funs       []FunDecl
	Pos        Pos
}

func (ImplBlock) isTopLevelDecl() {}

// FunDecl corresponds to: [export] fun name [<T,...>](params): Type { statements }
type FunDecl struct {
	Exported   bool
	Name       string
	TypeParams []string
	Params     []Param
	ReturnType TypeExpr
	Body       []Stmt
	Pos        Pos
}

func (FunDecl) isTopLevelDecl() {}

type Param struct {
	Name string
	Type TypeExpr
	Pos  Pos
}

// Type AST

type TypeExpr interface{ isType() }

type PrimitiveType struct {
	Name string
	Pos  Pos
}

func (PrimitiveType) isType() {}

type NamedType struct {
	Name string
	Pos  Pos
}

func (NamedType) isType() {}

type OptionalType struct {
	Inner TypeExpr
	Pos   Pos
}

func (OptionalType) isType() {}

type ArrayType struct {
	Inner TypeExpr
	Pos   Pos
}

func (ArrayType) isType() {}

type TupleType struct {
	Elems []TypeExpr
	Pos   Pos
}

func (TupleType) isType() {}

// Statements

type Stmt interface{ isStmt() }

type VarDecl struct {
	Name  string
	Ty    TypeExpr
	Value Expr
	Pos   Pos
}

func (VarDecl) isStmt() {}

type ExprStmt struct {
	E   Expr
	Pos Pos
}

func (ExprStmt) isStmt() {}

type HandleStmt struct {
	Expr    Expr
	OkName  string
	OkBody  []Stmt
	ErrName string
	ErrBody []Stmt
	Pos     Pos
}

func (HandleStmt) isStmt() {}

type ReturnStmt struct {
	Value Expr
	Pos   Pos
}

func (ReturnStmt) isStmt() {}

type IfStmt struct {
	Cond    Expr
	Then    []Stmt
	ElseIfs []IfStmt // only Cond+Then used for else-if chain
	Else    []Stmt
	Pos     Pos
}

func (IfStmt) isStmt() {}

type WhileStmt struct {
	Cond Expr
	Body []Stmt
	Pos  Pos
}

func (WhileStmt) isStmt() {}

type ForInStmt struct {
	Name     string
	Iterable Expr
	Body     []Stmt
	Pos      Pos
}

func (ForInStmt) isStmt() {}

type ForIStmt struct {
	Init VarDecl
	Cond Expr
	Post VarDecl
	Body []Stmt
	Pos  Pos
}

func (ForIStmt) isStmt() {}

// Additional expressions

type StringExpr struct {
	Value string
	Pos   Pos
}

func (StringExpr) isExpr() {}

type BoolExpr struct {
	Value bool
	Pos   Pos
}

func (BoolExpr) isExpr() {}

type MemberExpr struct {
	Object Expr
	Name   string
	Pos    Pos
}

func (MemberExpr) isExpr() {}

type CastExpr struct {
	Value Expr
	Ty    TypeExpr
	Pos   Pos
}

func (CastExpr) isExpr() {}

type ListExpr struct {
	Elems []Expr
	Pos   Pos
}

func (ListExpr) isExpr() {}

type TupleExpr struct {
	Elems []Expr
	Pos   Pos
}

func (TupleExpr) isExpr() {}

type UnaryExpr struct {
	Op    TokenType
	Right Expr
	Pos   Pos
}

func (UnaryExpr) isExpr() {}

type BinaryExpr struct {
	Left  Expr
	Op    TokenType
	Right Expr
	Pos   Pos
}

func (BinaryExpr) isExpr() {}
