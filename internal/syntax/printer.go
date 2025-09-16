package syntax

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Color styles for different AST elements
var (
	fileStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Bold(true) // Blue
	declStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))            // Green
	stmtStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))            // Yellow
	exprStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))             // Red
	identStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))            // Cyan
	literalStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("13"))            // Magenta
	keywordStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("5")).Bold(true)  // Purple
	fieldStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))             // Gray
)

// PrintFile returns a pretty-printed string representation of the AST
func PrintFile(file *File) string {
	if file == nil {
		return "<nil file>"
	}

	var builder strings.Builder
	builder.WriteString(fileStyle.Render("File") + " {\n")
	builder.WriteString(fmt.Sprintf("  %s: %s\n", fieldStyle.Render("Package"), identStyle.Render(file.Package)))
	builder.WriteString(fmt.Sprintf("  %s: [\n", fieldStyle.Render("Declarations")))

	for _, decl := range file.Decls {
		builder.WriteString(printDecl(decl, "    "))
	}

	builder.WriteString("  ]\n")
	builder.WriteString("}")

	return builder.String()
}

func printDecl(decl Decl, indent string) string {
	if decl == nil {
		return indent + "<nil decl>\n"
	}

	switch d := decl.(type) {
	case *FunDecl:
		return printFunDecl(d, indent)
	case *VarDecl:
		return printVarDecl(d, indent)
	default:
		return indent + fmt.Sprintf("UnknownDecl: %T\n", decl)
	}
}

func printFunDecl(decl *FunDecl, indent string) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("%s%s {\n", indent, declStyle.Render("FunDecl")))
	builder.WriteString(fmt.Sprintf("%s  %s: %s\n", indent, fieldStyle.Render("Name"), identStyle.Render(decl.Name)))
	builder.WriteString(fmt.Sprintf("%s  %s: %s\n", indent, fieldStyle.Render("Type"), identStyle.Render(decl.Type)))
	builder.WriteString(fmt.Sprintf("%s  %s: [\n", indent, fieldStyle.Render("Params")))

	for _, param := range decl.Params {
		builder.WriteString(fmt.Sprintf("%s    %s { %s: %s, %s: %s }\n",
			indent, keywordStyle.Render("Param"),
			fieldStyle.Render("Name"), identStyle.Render(param.Name),
			fieldStyle.Render("Type"), identStyle.Render(param.Type)))
	}

	builder.WriteString(fmt.Sprintf("%s  ]\n", indent))
	builder.WriteString(fmt.Sprintf("%s  %s: %s", indent, fieldStyle.Render("Body"), printStmt(decl.Body, indent+"  ")))
	builder.WriteString(fmt.Sprintf("%s}\n", indent))

	return builder.String()
}

func printVarDecl(decl *VarDecl, indent string) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("%sVarDecl {\n", indent))
	builder.WriteString(fmt.Sprintf("%s  Name: %s\n", indent, decl.Name))
	if decl.Type != "" {
		builder.WriteString(fmt.Sprintf("%s  Type: %s\n", indent, decl.Type))
	}
	builder.WriteString(fmt.Sprintf("%s  Init: %s", indent, printExpr(decl.Init, indent+"  ")))
	builder.WriteString(fmt.Sprintf("%s}\n", indent))

	return builder.String()
}

func printStmt(stmt Stmt, indent string) string {
	if stmt == nil {
		return indent + "<nil stmt>\n"
	}

	switch s := stmt.(type) {
	case *Block:
		return printBlock(s, indent)
	case *VarDecl:
		return printVarDecl(s, indent)
	case *AssignStmt:
		return printAssignStmt(s, indent)
	case *ExprStmt:
		return printExprStmt(s, indent)
	case *ReturnStmt:
		return printReturnStmt(s, indent)
	case *IfStmt:
		return printIfStmt(s, indent)
	case *WhileStmt:
		return printWhileStmt(s, indent)
	case *ForInStmt:
		return printForInStmt(s, indent)
	default:
		return indent + fmt.Sprintf("UnknownStmt: %T\n", stmt)
	}
}

func printBlock(stmt *Block, indent string) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("Block {\n"))

	for _, s := range stmt.Stmts {
		builder.WriteString(printStmt(s, indent+"  "))
	}

	builder.WriteString(fmt.Sprintf("%s}\n", indent))
	return builder.String()
}

func printAssignStmt(stmt *AssignStmt, indent string) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("%sAssignStmt {\n", indent))
	builder.WriteString(fmt.Sprintf("%s  Left: %s", indent, printExpr(stmt.Left, indent+"  ")))
	builder.WriteString(fmt.Sprintf("%s  Right: %s", indent, printExpr(stmt.Right, indent+"  ")))
	builder.WriteString(fmt.Sprintf("%s}\n", indent))

	return builder.String()
}

func printExprStmt(stmt *ExprStmt, indent string) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("%sExprStmt {\n", indent))
	builder.WriteString(fmt.Sprintf("%s  X: %s", indent, printExpr(stmt.X, indent+"  ")))
	builder.WriteString(fmt.Sprintf("%s}\n", indent))

	return builder.String()
}

func printReturnStmt(stmt *ReturnStmt, indent string) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("%sReturnStmt {\n", indent))
	builder.WriteString(fmt.Sprintf("%s  Result: %s", indent, printExpr(stmt.Result, indent+"  ")))
	builder.WriteString(fmt.Sprintf("%s}\n", indent))

	return builder.String()
}

func printIfStmt(stmt *IfStmt, indent string) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("%sIfStmt {\n", indent))
	builder.WriteString(fmt.Sprintf("%s  Cond: %s", indent, printExpr(stmt.Cond, indent+"  ")))
	builder.WriteString(fmt.Sprintf("%s  Body: %s", indent, printStmt(stmt.Body, indent+"  ")))

	if stmt.Else != nil {
		builder.WriteString(fmt.Sprintf("%s  Else: %s", indent, printStmt(stmt.Else, indent+"  ")))
	}

	builder.WriteString(fmt.Sprintf("%s}\n", indent))
	return builder.String()
}

func printWhileStmt(stmt *WhileStmt, indent string) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("%sWhileStmt {\n", indent))
	builder.WriteString(fmt.Sprintf("%s  Cond: %s", indent, printExpr(stmt.Cond, indent+"  ")))
	builder.WriteString(fmt.Sprintf("%s  Body: %s", indent, printStmt(stmt.Body, indent+"  ")))
	builder.WriteString(fmt.Sprintf("%s}\n", indent))

	return builder.String()
}

func printForInStmt(stmt *ForInStmt, indent string) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("%sForInStmt {\n", indent))
	builder.WriteString(fmt.Sprintf("%s  Var: %s\n", indent, stmt.Var))
	builder.WriteString(fmt.Sprintf("%s  Iter: %s", indent, printExpr(stmt.Iter, indent+"  ")))
	builder.WriteString(fmt.Sprintf("%s  Body: %s", indent, printStmt(stmt.Body, indent+"  ")))
	builder.WriteString(fmt.Sprintf("%s}\n", indent))

	return builder.String()
}

func printExpr(expr Expr, indent string) string {
	if expr == nil {
		return indent + "<nil expr>\n"
	}

	switch e := expr.(type) {
	case *BinaryExpr:
		return printBinaryExpr(e, indent)
	case *UnaryExpr:
		return printUnaryExpr(e, indent)
	case *CallExpr:
		return printCallExpr(e, indent)
	case *Ident:
		return printIdent(e, indent)
	case *BasicLit:
		return printBasicLit(e, indent)
	case *ArrayLit:
		return printArrayLit(e, indent)
	default:
		return indent + fmt.Sprintf("UnknownExpr: %T\n", expr)
	}
}

func printBinaryExpr(expr *BinaryExpr, indent string) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("BinaryExpr {\n"))
	builder.WriteString(fmt.Sprintf("%s  Op: %s\n", indent, expr.Op))
	builder.WriteString(fmt.Sprintf("%s  Left: %s", indent, printExpr(expr.Left, indent+"  ")))
	builder.WriteString(fmt.Sprintf("%s  Right: %s", indent, printExpr(expr.Right, indent+"  ")))
	builder.WriteString(fmt.Sprintf("%s}\n", indent))

	return builder.String()
}

func printUnaryExpr(expr *UnaryExpr, indent string) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("UnaryExpr {\n"))
	builder.WriteString(fmt.Sprintf("%s  Op: %s\n", indent, expr.Op))
	builder.WriteString(fmt.Sprintf("%s  X: %s", indent, printExpr(expr.X, indent+"  ")))
	builder.WriteString(fmt.Sprintf("%s}\n", indent))

	return builder.String()
}

func printCallExpr(expr *CallExpr, indent string) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("CallExpr {\n"))
	builder.WriteString(fmt.Sprintf("%s  Fun: %s", indent, printExpr(expr.Fun, indent+"  ")))
	builder.WriteString(fmt.Sprintf("%s  Args: [\n", indent))

	for _, arg := range expr.Args {
		builder.WriteString(fmt.Sprintf("%s    %s", indent, printExpr(arg, indent+"    ")))
	}

	builder.WriteString(fmt.Sprintf("%s  ]\n", indent))
	builder.WriteString(fmt.Sprintf("%s}\n", indent))

	return builder.String()
}

func printIdent(expr *Ident, indent string) string {
	return fmt.Sprintf("%s { %s: %s }\n",
		exprStyle.Render("Ident"),
		fieldStyle.Render("Name"),
		identStyle.Render(expr.Name))
}

func printBasicLit(expr *BasicLit, indent string) string {
	return fmt.Sprintf("%s { %s: %s, %s: %s }\n",
		exprStyle.Render("BasicLit"),
		fieldStyle.Render("Kind"),
		keywordStyle.Render(string(expr.Kind)),
		fieldStyle.Render("Value"),
		literalStyle.Render(expr.Value))
}

func printArrayLit(expr *ArrayLit, indent string) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("ArrayLit {\n"))
	builder.WriteString(fmt.Sprintf("%s  Elements: [\n", indent))

	for _, elem := range expr.Elements {
		builder.WriteString(fmt.Sprintf("%s    %s", indent, printExpr(elem, indent+"    ")))
	}

	builder.WriteString(fmt.Sprintf("%s  ]\n", indent))
	builder.WriteString(fmt.Sprintf("%s}\n", indent))

	return builder.String()
}
