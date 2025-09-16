package syntax

import "jmpeax.com/guayavita/gvc/internal/diag"

// Core AST interfaces
type Node interface {
	Pos() diag.Position
}

type Expr interface {
	Node
	exprNode()
}

type Stmt interface {
	Node
	stmtNode()
}

type Decl interface {
	Node
	declNode()
}

// File represents a complete source file
type File struct {
	Package string
	Decls   []Decl
	Pos_    diag.Position
}

func (f *File) Pos() diag.Position { return f.Pos_ }

// Declarations
type FunDecl struct {
	Name   string
	Params []Param
	Type   string
	Body   *Block
	Pos_   diag.Position
}

func (d *FunDecl) Pos() diag.Position { return d.Pos_ }
func (d *FunDecl) declNode()          {}

type VarDecl struct {
	Name string
	Type string // optional, empty if not specified
	Init Expr
	Pos_ diag.Position
}

func (d *VarDecl) Pos() diag.Position { return d.Pos_ }
func (d *VarDecl) declNode()          {}
func (d *VarDecl) stmtNode()          {}

type Param struct {
	Name string
	Type string
	Pos_ diag.Position
}

func (p *Param) Pos() diag.Position { return p.Pos_ }

// Statements
type Block struct {
	Stmts []Stmt
	Pos_  diag.Position
}

func (s *Block) Pos() diag.Position { return s.Pos_ }
func (s *Block) stmtNode()          {}

type AssignStmt struct {
	Left  Expr
	Right Expr
	Pos_  diag.Position
}

func (s *AssignStmt) Pos() diag.Position { return s.Pos_ }
func (s *AssignStmt) stmtNode()          {}

type ExprStmt struct {
	X    Expr
	Pos_ diag.Position
}

func (s *ExprStmt) Pos() diag.Position { return s.Pos_ }
func (s *ExprStmt) stmtNode()          {}

type ReturnStmt struct {
	Result Expr
	Pos_   diag.Position
}

func (s *ReturnStmt) Pos() diag.Position { return s.Pos_ }
func (s *ReturnStmt) stmtNode()          {}

type IfStmt struct {
	Cond Expr
	Body *Block
	Else Stmt // can be another IfStmt for "else if" or Block for "else"
	Pos_ diag.Position
}

func (s *IfStmt) Pos() diag.Position { return s.Pos_ }
func (s *IfStmt) stmtNode()          {}

type WhileStmt struct {
	Cond Expr
	Body *Block
	Pos_ diag.Position
}

func (s *WhileStmt) Pos() diag.Position { return s.Pos_ }
func (s *WhileStmt) stmtNode()          {}

type ForInStmt struct {
	Var  string
	Iter Expr
	Body *Block
	Pos_ diag.Position
}

func (s *ForInStmt) Pos() diag.Position { return s.Pos_ }
func (s *ForInStmt) stmtNode()          {}

// Expressions
type BinaryExpr struct {
	Left  Expr
	Op    string
	Right Expr
	Pos_  diag.Position
}

func (e *BinaryExpr) Pos() diag.Position { return e.Pos_ }
func (e *BinaryExpr) exprNode()          {}

type UnaryExpr struct {
	Op   string
	X    Expr
	Pos_ diag.Position
}

func (e *UnaryExpr) Pos() diag.Position { return e.Pos_ }
func (e *UnaryExpr) exprNode()          {}

type CallExpr struct {
	Fun  Expr
	Args []Expr
	Pos_ diag.Position
}

func (e *CallExpr) Pos() diag.Position { return e.Pos_ }
func (e *CallExpr) exprNode()          {}

type Ident struct {
	Name string
	Pos_ diag.Position
}

func (e *Ident) Pos() diag.Position { return e.Pos_ }
func (e *Ident) exprNode()          {}

type BasicLit struct {
	Kind  string // "INT", "FLOAT", "STRING", "BOOL", "NONE"
	Value string
	Pos_  diag.Position
}

func (e *BasicLit) Pos() diag.Position { return e.Pos_ }
func (e *BasicLit) exprNode()          {}

type ArrayLit struct {
	Elements []Expr
	Pos_     diag.Position
}

func (e *ArrayLit) Pos() diag.Position { return e.Pos_ }
func (e *ArrayLit) exprNode()          {}
