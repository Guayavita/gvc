package syntax

import (
	"strconv"
	"strings"

	"jmpeax.com/guayavita/gvc/internal/diag"
)

type TokenKind string

const (
	// Special tokens
	ILLEGAL TokenKind = "ILLEGAL"
	EOF     TokenKind = "EOF"

	// Identifiers and literals
	IDENT  TokenKind = "IDENT"
	INT    TokenKind = "INT"
	FLOAT  TokenKind = "FLOAT"
	STRING TokenKind = "STRING"
	TRUE   TokenKind = "TRUE"
	FALSE  TokenKind = "FALSE"
	NONE   TokenKind = "NONE"

	// Keywords
	PACKAGE TokenKind = "PACKAGE"
	IMPORT  TokenKind = "IMPORT"
	DEF     TokenKind = "DEF"
	FUN     TokenKind = "FUN"
	TYPE    TokenKind = "TYPE"
	EXPORT  TokenKind = "EXPORT"
	RETURN  TokenKind = "RETURN"
	IF      TokenKind = "IF"
	ELSE    TokenKind = "ELSE"
	WHILE   TokenKind = "WHILE"
	FOR     TokenKind = "FOR"
	IN      TokenKind = "IN"
	HANDLE  TokenKind = "HANDLE"
	OK      TokenKind = "OK"
	ERR     TokenKind = "ERR"
	STRUCT  TokenKind = "STRUCT"
	ENUM    TokenKind = "ENUM"
	IMPL    TokenKind = "IMPL"
	AS      TokenKind = "AS"

	// Operators
	ASSIGN TokenKind = "="
	EQ     TokenKind = "=="
	NE     TokenKind = "!="
	LT     TokenKind = "<"
	LE     TokenKind = "<="
	GT     TokenKind = ">"
	GE     TokenKind = ">="
	PLUS   TokenKind = "+"
	MINUS  TokenKind = "-"
	MUL    TokenKind = "*"
	DIV    TokenKind = "/"
	MOD    TokenKind = "%"
	AND    TokenKind = "&&"
	OR     TokenKind = "||"
	NOT    TokenKind = "!"

	// Punctuation
	COMMA     TokenKind = ","
	SEMICOLON TokenKind = ";"
	COLON     TokenKind = ":"
	DOT       TokenKind = "."
	ARROW     TokenKind = "->"
	QUESTION  TokenKind = "?"

	// Delimiters
	LPAREN   TokenKind = "("
	RPAREN   TokenKind = ")"
	LBRACE   TokenKind = "{"
	RBRACE   TokenKind = "}"
	LBRACKET TokenKind = "["
	RBRACKET TokenKind = "]"
)

type Token struct {
	Kind  TokenKind
	Value string
	Pos   diag.Position
}

var keywords = map[string]TokenKind{
	"package": PACKAGE,
	"import":  IMPORT,
	"def":     DEF,
	"fun":     FUN,
	"type":    TYPE,
	"export":  EXPORT,
	"return":  RETURN,
	"if":      IF,
	"else":    ELSE,
	"while":   WHILE,
	"for":     FOR,
	"in":      IN,
	"handle":  HANDLE,
	"Ok":      OK,
	"Err":     ERR,
	"struct":  STRUCT,
	"enum":    ENUM,
	"impl":    IMPL,
	"as":      AS,
	"true":    TRUE,
	"false":   FALSE,
	"none":    NONE,
}

type Lexer struct {
	input    string
	filename string
	pos      int  // current position in input (points to current char)
	readPos  int  // current reading position in input (after current char)
	ch       byte // current char under examination
	line     int
	column   int
}

func NewLexer(input, filename string) *Lexer {
	l := &Lexer{
		input:    input,
		filename: filename,
		line:     1,
		column:   0,
	}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPos >= len(l.input) {
		l.ch = 0 // ASCII NUL character represents "EOF"
	} else {
		l.ch = l.input[l.readPos]
	}
	l.pos = l.readPos
	l.readPos++
	if l.ch == '\n' {
		l.line++
		l.column = 0
	} else {
		l.column++
	}
}

func (l *Lexer) peekChar() byte {
	if l.readPos >= len(l.input) {
		return 0
	}
	return l.input[l.readPos]
}

func (l *Lexer) currentPos() diag.Position {
	return diag.Position{
		File:   l.filename,
		Line:   l.line,
		Column: l.column,
	}
}

func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()
	l.skipComments()

	tok.Pos = l.currentPos()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{Kind: EQ, Value: "==", Pos: tok.Pos}
		} else {
			tok = Token{Kind: ASSIGN, Value: string(l.ch), Pos: tok.Pos}
		}
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{Kind: NE, Value: "!=", Pos: tok.Pos}
		} else {
			tok = Token{Kind: NOT, Value: string(l.ch), Pos: tok.Pos}
		}
	case '<':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{Kind: LE, Value: "<=", Pos: tok.Pos}
		} else {
			tok = Token{Kind: LT, Value: string(l.ch), Pos: tok.Pos}
		}
	case '>':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{Kind: GE, Value: ">=", Pos: tok.Pos}
		} else {
			tok = Token{Kind: GT, Value: string(l.ch), Pos: tok.Pos}
		}
	case '&':
		if l.peekChar() == '&' {
			l.readChar()
			tok = Token{Kind: AND, Value: "&&", Pos: tok.Pos}
		} else {
			tok = Token{Kind: ILLEGAL, Value: string(l.ch), Pos: tok.Pos}
		}
	case '|':
		if l.peekChar() == '|' {
			l.readChar()
			tok = Token{Kind: OR, Value: "||", Pos: tok.Pos}
		} else {
			tok = Token{Kind: ILLEGAL, Value: string(l.ch), Pos: tok.Pos}
		}
	case '-':
		if l.peekChar() == '>' {
			l.readChar()
			tok = Token{Kind: ARROW, Value: "->", Pos: tok.Pos}
		} else {
			tok = Token{Kind: MINUS, Value: string(l.ch), Pos: tok.Pos}
		}
	case '+':
		tok = Token{Kind: PLUS, Value: string(l.ch), Pos: tok.Pos}
	case '*':
		tok = Token{Kind: MUL, Value: string(l.ch), Pos: tok.Pos}
	case '/':
		tok = Token{Kind: DIV, Value: string(l.ch), Pos: tok.Pos}
	case '%':
		tok = Token{Kind: MOD, Value: string(l.ch), Pos: tok.Pos}
	case ',':
		tok = Token{Kind: COMMA, Value: string(l.ch), Pos: tok.Pos}
	case ';':
		tok = Token{Kind: SEMICOLON, Value: string(l.ch), Pos: tok.Pos}
	case ':':
		tok = Token{Kind: COLON, Value: string(l.ch), Pos: tok.Pos}
	case '.':
		tok = Token{Kind: DOT, Value: string(l.ch), Pos: tok.Pos}
	case '?':
		tok = Token{Kind: QUESTION, Value: string(l.ch), Pos: tok.Pos}
	case '(':
		tok = Token{Kind: LPAREN, Value: string(l.ch), Pos: tok.Pos}
	case ')':
		tok = Token{Kind: RPAREN, Value: string(l.ch), Pos: tok.Pos}
	case '{':
		tok = Token{Kind: LBRACE, Value: string(l.ch), Pos: tok.Pos}
	case '}':
		tok = Token{Kind: RBRACE, Value: string(l.ch), Pos: tok.Pos}
	case '[':
		tok = Token{Kind: LBRACKET, Value: string(l.ch), Pos: tok.Pos}
	case ']':
		tok = Token{Kind: RBRACKET, Value: string(l.ch), Pos: tok.Pos}
	case '"':
		tok.Value = l.readString()
		tok.Kind = STRING
		return tok // readString() advances position
	case 0:
		tok = Token{Kind: EOF, Value: "", Pos: tok.Pos}
	default:
		if isLetter(l.ch) {
			tok.Value = l.readIdentifier()
			tok.Kind = lookupIdent(tok.Value)
			return tok // readIdentifier() advances position
		} else if isDigit(l.ch) {
			tok.Value = l.readNumber()
			if strings.Contains(tok.Value, ".") {
				tok.Kind = FLOAT
			} else {
				tok.Kind = INT
			}
			return tok // readNumber() advances position
		} else {
			tok = Token{Kind: ILLEGAL, Value: string(l.ch), Pos: tok.Pos}
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) skipComments() {
	if l.ch == '/' {
		if l.peekChar() == '/' {
			// Line comment
			for l.ch != '\n' && l.ch != 0 {
				l.readChar()
			}
			l.skipWhitespace()
		} else if l.peekChar() == '*' {
			// Block comment
			l.readChar() // consume '/'
			l.readChar() // consume '*'
			for {
				if l.ch == '*' && l.peekChar() == '/' {
					l.readChar() // consume '*'
					l.readChar() // consume '/'
					break
				}
				if l.ch == 0 {
					break // EOF in comment
				}
				l.readChar()
			}
			l.skipWhitespace()
		}
	}
}

func (l *Lexer) readIdentifier() string {
	position := l.pos
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' {
		l.readChar()
	}
	return l.input[position:l.pos]
}

func (l *Lexer) readNumber() string {
	position := l.pos
	for isDigit(l.ch) {
		l.readChar()
	}

	// Handle float
	if l.ch == '.' && isDigit(l.peekChar()) {
		l.readChar() // consume '.'
		for isDigit(l.ch) {
			l.readChar()
		}
	}

	return l.input[position:l.pos]
}

func (l *Lexer) readString() string {
	l.readChar() // move past opening quote

	var result strings.Builder
	for l.ch != '"' && l.ch != 0 {
		if l.ch == '\\' {
			l.readChar()
			switch l.ch {
			case 'n':
				result.WriteByte('\n')
			case 't':
				result.WriteByte('\t')
			case 'r':
				result.WriteByte('\r')
			case '"':
				result.WriteByte('"')
			case '\\':
				result.WriteByte('\\')
			case 'u':
				// Unicode escape \uXXXX
				l.readChar()
				hex := l.readHex(4)
				if code, err := strconv.ParseUint(hex, 16, 16); err == nil {
					result.WriteRune(rune(code))
				}
				continue // readHex already advanced position
			case 'U':
				// Unicode escape \UXXXXXXXX
				l.readChar()
				hex := l.readHex(8)
				if code, err := strconv.ParseUint(hex, 16, 32); err == nil {
					result.WriteRune(rune(code))
				}
				continue // readHex already advanced position
			case 'x':
				// Hex escape \xXX
				l.readChar()
				hex := l.readHex(2)
				if code, err := strconv.ParseUint(hex, 16, 8); err == nil {
					result.WriteByte(byte(code))
				}
				continue // readHex already advanced position
			default:
				result.WriteByte(l.ch)
			}
		} else {
			result.WriteByte(l.ch)
		}
		l.readChar()
	}

	if l.ch == '"' {
		l.readChar() // consume closing quote
	}

	return result.String()
}

func (l *Lexer) readHex(count int) string {
	position := l.pos
	for i := 0; i < count && isHexDigit(l.ch); i++ {
		l.readChar()
	}
	return l.input[position:l.pos]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isHexDigit(ch byte) bool {
	return isDigit(ch) || 'a' <= ch && ch <= 'f' || 'A' <= ch && ch <= 'F'
}

func lookupIdent(ident string) TokenKind {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
