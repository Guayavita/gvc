package syntax

import (
	"jmpeax.com/guayavita/gvc/internal/diag"
)

type Parser struct {
	lexer       *Lexer
	curToken    Token
	peekToken   Token
	diagnostics []diag.Diagnostic
}

// ParseFile parses a Guayavita source file and returns the AST and any diagnostics
func ParseFile(filename, source string) (*File, []diag.Diagnostic) {
	lexer := NewLexer(source, filename)
	parser := &Parser{
		lexer:       lexer,
		diagnostics: []diag.Diagnostic{},
	}

	// Read two tokens to initialize current and peek
	parser.nextToken()
	parser.nextToken()

	file := parser.parseFile()
	return file, parser.diagnostics
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) isTypeKeyword(kind TokenKind) bool {
	switch kind {
	case NONE:
		return true
	default:
		// Check if it's a primitive type by looking at the token value
		switch p.curToken.Value {
		case "bool", "i8", "i32", "i64", "u8", "u16", "u32", "u64",
			"f32", "f64", "byte", "string":
			return true
		default:
			return false
		}
	}
}

func (p *Parser) error(msg string) {
	diagnostic := diag.Diagnostic{
		Severity: diag.Error,
		Message:  msg,
		Span: diag.Span{
			Start: p.curToken.Pos,
			End:   p.curToken.Pos,
		},
	}
	p.diagnostics = append(p.diagnostics, diagnostic)
}

func (p *Parser) expectToken(expected TokenKind) bool {
	if p.curToken.Kind != expected {
		p.error("expected " + string(expected) + ", got " + string(p.curToken.Kind))
		return false
	}
	return true
}

func (p *Parser) parseFile() *File {
	file := &File{
		Pos_:  p.curToken.Pos,
		Decls: []Decl{},
	}

	// Parse package declaration
	if p.curToken.Kind == PACKAGE {
		p.nextToken() // consume 'package'
		if p.expectToken(IDENT) {
			file.Package = p.curToken.Value
			p.nextToken()
		}
	}

	// Parse declarations
	for p.curToken.Kind != EOF {
		decl := p.parseDecl()
		if decl != nil {
			file.Decls = append(file.Decls, decl)
		}
	}

	return file
}

func (p *Parser) parseDecl() Decl {
	switch p.curToken.Kind {
	case DEF:
		return p.parseVarDecl()
	case FUN:
		return p.parseFunDecl()
	default:
		p.error("expected declaration, got " + string(p.curToken.Kind))
		p.nextToken() // skip invalid token
		return nil
	}
}

func (p *Parser) parseVarDecl() *VarDecl {
	pos := p.curToken.Pos
	p.nextToken() // consume 'def'

	if !p.expectToken(IDENT) {
		return nil
	}

	name := p.curToken.Value
	p.nextToken()

	var typeName string
	if p.curToken.Kind == COLON {
		p.nextToken() // consume ':'
		if p.curToken.Kind == IDENT || p.isTypeKeyword(p.curToken.Kind) {
			typeName = p.curToken.Value
			p.nextToken()
		} else {
			p.error("expected type identifier")
		}
	}

	if !p.expectToken(ASSIGN) {
		return nil
	}
	p.nextToken() // consume '='

	init := p.parseExpr()

	return &VarDecl{
		Name: name,
		Type: typeName,
		Init: init,
		Pos_: pos,
	}
}

func (p *Parser) parseFunDecl() *FunDecl {
	pos := p.curToken.Pos
	p.nextToken() // consume 'fun'

	if !p.expectToken(IDENT) {
		return nil
	}

	name := p.curToken.Value
	p.nextToken()

	if !p.expectToken(LPAREN) {
		return nil
	}
	p.nextToken() // consume '('

	params := []Param{}
	for p.curToken.Kind != RPAREN && p.curToken.Kind != EOF {
		param := p.parseParam()
		if param != nil {
			params = append(params, *param)
		}

		if p.curToken.Kind == COMMA {
			p.nextToken()
		} else if p.curToken.Kind != RPAREN {
			p.error("expected ',' or ')' in parameter list")
			break
		}
	}

	if !p.expectToken(RPAREN) {
		return nil
	}
	p.nextToken() // consume ')'

	if !p.expectToken(COLON) {
		return nil
	}
	p.nextToken() // consume ':'

	var returnType string
	if p.curToken.Kind == IDENT || p.curToken.Kind == NONE {
		returnType = p.curToken.Value
		p.nextToken()
	} else {
		p.error("expected type identifier")
	}

	body := p.parseBlock()

	return &FunDecl{
		Name:   name,
		Params: params,
		Type:   returnType,
		Body:   body,
		Pos_:   pos,
	}
}

func (p *Parser) parseParam() *Param {
	if !p.expectToken(IDENT) {
		return nil
	}

	pos := p.curToken.Pos
	name := p.curToken.Value
	p.nextToken()

	if !p.expectToken(COLON) {
		return nil
	}
	p.nextToken() // consume ':'

	if p.curToken.Kind == IDENT || p.isTypeKeyword(p.curToken.Kind) {
		typeName := p.curToken.Value
		p.nextToken()

		return &Param{
			Name: name,
			Type: typeName,
			Pos_: pos,
		}
	} else {
		p.error("expected type identifier")
		return nil
	}
}

func (p *Parser) parseBlock() *Block {
	pos := p.curToken.Pos

	if !p.expectToken(LBRACE) {
		return nil
	}
	p.nextToken() // consume '{'

	stmts := []Stmt{}
	for p.curToken.Kind != RBRACE && p.curToken.Kind != EOF {
		stmt := p.parseStmt()
		if stmt != nil {
			stmts = append(stmts, stmt)
		}
	}

	if !p.expectToken(RBRACE) {
		return nil
	}
	p.nextToken() // consume '}'

	return &Block{
		Stmts: stmts,
		Pos_:  pos,
	}
}

func (p *Parser) parseStmt() Stmt {
	switch p.curToken.Kind {
	case DEF:
		return p.parseVarDecl()
	case RETURN:
		return p.parseReturnStmt()
	case IF:
		return p.parseIfStmt()
	case WHILE:
		return p.parseWhileStmt()
	case FOR:
		return p.parseForStmt()
	default:
		// Expression statement
		expr := p.parseExpr()
		return &ExprStmt{
			X:    expr,
			Pos_: expr.Pos(),
		}
	}
}

func (p *Parser) parseReturnStmt() *ReturnStmt {
	pos := p.curToken.Pos
	p.nextToken() // consume 'return'

	result := p.parseExpr()

	return &ReturnStmt{
		Result: result,
		Pos_:   pos,
	}
}

func (p *Parser) parseIfStmt() *IfStmt {
	pos := p.curToken.Pos
	p.nextToken() // consume 'if'

	cond := p.parseExpr()
	body := p.parseBlock()

	var elseStmt Stmt
	if p.curToken.Kind == ELSE {
		p.nextToken() // consume 'else'
		if p.curToken.Kind == IF {
			elseStmt = p.parseIfStmt()
		} else {
			elseStmt = p.parseBlock()
		}
	}

	return &IfStmt{
		Cond: cond,
		Body: body,
		Else: elseStmt,
		Pos_: pos,
	}
}

func (p *Parser) parseWhileStmt() *WhileStmt {
	pos := p.curToken.Pos
	p.nextToken() // consume 'while'

	cond := p.parseExpr()
	body := p.parseBlock()

	return &WhileStmt{
		Cond: cond,
		Body: body,
		Pos_: pos,
	}
}

func (p *Parser) parseForStmt() Stmt {
	pos := p.curToken.Pos
	p.nextToken() // consume 'for'

	// Check for 'for def var in expr' pattern
	if p.curToken.Kind == DEF {
		p.nextToken() // consume 'def'

		if !p.expectToken(IDENT) {
			return nil
		}

		varName := p.curToken.Value
		p.nextToken()

		if !p.expectToken(IN) {
			return nil
		}
		p.nextToken() // consume 'in'

		iter := p.parseExpr()
		body := p.parseBlock()

		return &ForInStmt{
			Var:  varName,
			Iter: iter,
			Body: body,
			Pos_: pos,
		}
	}

	// For now, just skip other for loop forms
	p.error("unsupported for loop syntax")
	return nil
}

func (p *Parser) parseExpr() Expr {
	return p.parseOrExpr()
}

func (p *Parser) parseOrExpr() Expr {
	left := p.parseAndExpr()

	for p.curToken.Kind == OR {
		op := p.curToken.Value
		pos := p.curToken.Pos
		p.nextToken()
		right := p.parseAndExpr()
		left = &BinaryExpr{
			Left:  left,
			Op:    op,
			Right: right,
			Pos_:  pos,
		}
	}

	return left
}

func (p *Parser) parseAndExpr() Expr {
	left := p.parseCmpExpr()

	for p.curToken.Kind == AND {
		op := p.curToken.Value
		pos := p.curToken.Pos
		p.nextToken()
		right := p.parseCmpExpr()
		left = &BinaryExpr{
			Left:  left,
			Op:    op,
			Right: right,
			Pos_:  pos,
		}
	}

	return left
}

func (p *Parser) parseCmpExpr() Expr {
	left := p.parseAddExpr()

	for p.curToken.Kind == EQ || p.curToken.Kind == NE || p.curToken.Kind == LT ||
		p.curToken.Kind == LE || p.curToken.Kind == GT || p.curToken.Kind == GE {
		op := p.curToken.Value
		pos := p.curToken.Pos
		p.nextToken()
		right := p.parseAddExpr()
		left = &BinaryExpr{
			Left:  left,
			Op:    op,
			Right: right,
			Pos_:  pos,
		}
	}

	return left
}

func (p *Parser) parseAddExpr() Expr {
	left := p.parseMulExpr()

	for p.curToken.Kind == PLUS || p.curToken.Kind == MINUS {
		op := p.curToken.Value
		pos := p.curToken.Pos
		p.nextToken()
		right := p.parseMulExpr()
		left = &BinaryExpr{
			Left:  left,
			Op:    op,
			Right: right,
			Pos_:  pos,
		}
	}

	return left
}

func (p *Parser) parseMulExpr() Expr {
	left := p.parseUnaryExpr()

	for p.curToken.Kind == MUL || p.curToken.Kind == DIV || p.curToken.Kind == MOD {
		op := p.curToken.Value
		pos := p.curToken.Pos
		p.nextToken()
		right := p.parseUnaryExpr()
		left = &BinaryExpr{
			Left:  left,
			Op:    op,
			Right: right,
			Pos_:  pos,
		}
	}

	return left
}

func (p *Parser) parseUnaryExpr() Expr {
	if p.curToken.Kind == NOT || p.curToken.Kind == MINUS || p.curToken.Kind == PLUS {
		op := p.curToken.Value
		pos := p.curToken.Pos
		p.nextToken()
		expr := p.parseUnaryExpr()
		return &UnaryExpr{
			Op:   op,
			X:    expr,
			Pos_: pos,
		}
	}

	return p.parsePostfixExpr()
}

func (p *Parser) parsePostfixExpr() Expr {
	left := p.parsePrimary()

	for {
		switch p.curToken.Kind {
		case LPAREN:
			// Function call
			p.nextToken() // consume '('
			args := []Expr{}

			for p.curToken.Kind != RPAREN && p.curToken.Kind != EOF {
				arg := p.parseExpr()
				args = append(args, arg)

				if p.curToken.Kind == COMMA {
					p.nextToken()
				} else if p.curToken.Kind != RPAREN {
					p.error("expected ',' or ')' in argument list")
					break
				}
			}

			if !p.expectToken(RPAREN) {
				return left
			}
			p.nextToken() // consume ')'

			left = &CallExpr{
				Fun:  left,
				Args: args,
				Pos_: left.Pos(),
			}
		default:
			return left
		}
	}
}

func (p *Parser) parsePrimary() Expr {
	switch p.curToken.Kind {
	case IDENT:
		ident := &Ident{
			Name: p.curToken.Value,
			Pos_: p.curToken.Pos,
		}
		p.nextToken()
		return ident

	case INT, FLOAT, STRING:
		lit := &BasicLit{
			Kind:  string(p.curToken.Kind),
			Value: p.curToken.Value,
			Pos_:  p.curToken.Pos,
		}
		p.nextToken()
		return lit

	case TRUE, FALSE:
		lit := &BasicLit{
			Kind:  "BOOL",
			Value: p.curToken.Value,
			Pos_:  p.curToken.Pos,
		}
		p.nextToken()
		return lit

	case NONE:
		lit := &BasicLit{
			Kind:  "NONE",
			Value: p.curToken.Value,
			Pos_:  p.curToken.Pos,
		}
		p.nextToken()
		return lit

	case LBRACKET:
		return p.parseArrayLit()

	case LPAREN:
		p.nextToken() // consume '('
		expr := p.parseExpr()
		if !p.expectToken(RPAREN) {
			return expr
		}
		p.nextToken() // consume ')'
		return expr

	default:
		p.error("unexpected token in expression: " + string(p.curToken.Kind))
		p.nextToken() // skip invalid token
		return &BasicLit{Kind: "INVALID", Value: "", Pos_: p.curToken.Pos}
	}
}

func (p *Parser) parseArrayLit() *ArrayLit {
	pos := p.curToken.Pos
	p.nextToken() // consume '['

	elements := []Expr{}
	for p.curToken.Kind != RBRACKET && p.curToken.Kind != EOF {
		elem := p.parseExpr()
		elements = append(elements, elem)

		if p.curToken.Kind == COMMA {
			p.nextToken()
		} else if p.curToken.Kind != RBRACKET {
			p.error("expected ',' or ']' in array literal")
			break
		}
	}

	if !p.expectToken(RBRACKET) {
		return nil
	}
	p.nextToken() // consume ']'

	return &ArrayLit{
		Elements: elements,
		Pos_:     pos,
	}
}
