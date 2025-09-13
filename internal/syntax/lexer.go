package syntax

import (
	"bufio"
	"io"
	"strings"
	"unicode"
)

// TokenType represents the kinds of tokens our lexer can produce.
type TokenType int

const (
	// Special
	EOF TokenType = iota
	ILLEGAL

	// Keywords
	PACKAGE
	IMPORT
	EXPORT
	DEF
	TYPE
	STRUCT
	ENUM
	IMPL
	FUN
	RETURN
	IF
	ELSE
	WHILE
	FOR
	IN
	HANDLE
	OK
	ERR
	AS
	TRUE
	FALSE

	// Identifiers
	IDENT

	// Literals
	NUMBER
	STRING

	// Symbols / punctuation
	ASSIGN    // '='
	PIPE      // '|'
	LPAREN    // '('
	RPAREN    // ')'
	LBRACE    // '{'
	RBRACE    // '}'
	LBRACKET  // '['
	RBRACKET  // ']'
	LANGLE    // '<'
	RANGLE    // '>'
	COLON     // ':'
	SEMICOLON // ';'
	COMMA     // ','
	DOT       // '.'
	QUESTION  // '?'
	BANG      // '!'
	PLUS      // '+'
	MINUS     // '-'
	STAR      // '*'
	SLASH     // '/'
	EQEQ      // '=='
	NEQ       // '!='
	LTE       // '<='
	GTE       // '>='
	OROR      // '||'
	ANDAND    // '&&'
	PERCENT   // '%'
)

func (t TokenType) String() string {
	switch t {
	case EOF:
		return "EOF"
	case ILLEGAL:
		return "ILLEGAL"
	case PACKAGE:
		return "PACKAGE"
	case IMPORT:
		return "IMPORT"
	case EXPORT:
		return "EXPORT"
	case DEF:
		return "DEF"
	case TYPE:
		return "TYPE"
	case STRUCT:
		return "STRUCT"
	case ENUM:
		return "ENUM"
	case IMPL:
		return "IMPL"
	case FUN:
		return "FUN"
	case RETURN:
		return "RETURN"
	case IF:
		return "IF"
	case ELSE:
		return "ELSE"
	case WHILE:
		return "WHILE"
	case FOR:
		return "FOR"
	case IN:
		return "IN"
	case HANDLE:
		return "HANDLE"
	case OK:
		return "OK"
	case ERR:
		return "ERR"
	case AS:
		return "AS"
	case TRUE:
		return "TRUE"
	case FALSE:
		return "FALSE"
	case IDENT:
		return "IDENT"
	case NUMBER:
		return "NUMBER"
	case STRING:
		return "STRING"
	case ASSIGN:
		return "ASSIGN"
	case PIPE:
		return "PIPE"
	case LPAREN:
		return "LPAREN"
	case RPAREN:
		return "RPAREN"
	case LBRACE:
		return "LBRACE"
	case RBRACE:
		return "RBRACE"
	case LBRACKET:
		return "LBRACKET"
	case RBRACKET:
		return "RBRACKET"
	case LANGLE:
		return "LANGLE"
	case RANGLE:
		return "RANGLE"
	case COLON:
		return "COLON"
	case SEMICOLON:
		return "SEMICOLON"
	case COMMA:
		return "COMMA"
	case DOT:
		return "DOT"
	case QUESTION:
		return "QUESTION"
	case BANG:
		return "BANG"
	case PLUS:
		return "PLUS"
	case MINUS:
		return "MINUS"
	case STAR:
		return "STAR"
	case SLASH:
		return "SLASH"
	case EQEQ:
		return "EQEQ"
	case NEQ:
		return "NEQ"
	case LTE:
		return "LTE"
	case GTE:
		return "GTE"
	case OROR:
		return "OROR"
	case ANDAND:
		return "ANDAND"
	case PERCENT:
		return "PERCENT"
	default:
		// Minimal itoa to avoid importing fmt
		return "TokenType(" + itoa(int(t)) + ")"
	}
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	neg := false
	if i < 0 {
		neg = true
		i = -i
	}
	var buf [20]byte
	b := len(buf)
	for i > 0 {
		b--
		buf[b] = byte('0' + i%10)
		i /= 10
	}
	if neg {
		b--
		buf[b] = '-'
	}
	return string(buf[b:])
}

type Token struct {
	Type TokenType
	Lit  string
	Pos  Pos
}

// Lexer implements a simple rune-based lexer for the grammar: "package NAME".
type Lexer struct {
	r      *bufio.Reader
	offset int
	line   int
	col    int

	peekedRune rune
	peekedSize int
	hasPeeked  bool
}

func NewLexerFromString(s string) *Lexer {
	return NewLexer(strings.NewReader(s))
}

func NewLexer(r io.Reader) *Lexer {
	return &Lexer{
		r:    bufio.NewReader(r),
		line: 1,
		col:  0,
	}
}

func (lx *Lexer) readRune() (r rune, size int, err error) {
	if lx.hasPeeked {
		r, size = lx.peekedRune, lx.peekedSize
		lx.hasPeeked = false
		return r, size, nil
	}
	r, size, err = lx.r.ReadRune()
	if err != nil {
		return 0, 0, err
	}
	lx.offset += size
	if r == '\n' {
		lx.line++
		lx.col = 0
	} else {
		lx.col++
	}
	return r, size, nil
}

func (lx *Lexer) unreadRune(r rune, size int) {
	if lx.hasPeeked {
		panic("internal lexer error: double unread")
	}
	lx.peekedRune = r
	lx.peekedSize = size
	lx.hasPeeked = true

	lx.offset -= size
	if r == '\n' {
		lx.line--
		lx.col = 0
	} else if lx.col > 0 {
		lx.col--
	}
}

func (lx *Lexer) skipWhitespaceAndComments() error {
	for {
		r, size, err := lx.readRune()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		// Skip whitespace only; comments are handled in NextToken to avoid double unread
		if unicode.IsSpace(r) {
			continue
		}
		// Non-whitespace: unread it and return for tokenization
		lx.unreadRune(r, size)
		return nil
	}
}

func (lx *Lexer) NextToken() (Token, error) {
	_ = lx.skipWhitespaceAndComments() // best effort; errors surface while reading runes

	startPos := Pos{Offset: lx.offset, Line: lx.line, Col: lx.col + 1}

	r, _, err := lx.readRune()
	if err == io.EOF {
		return Token{Type: EOF, Pos: startPos}, nil
	}
	if err != nil {
		return Token{}, err
	}

	// Handle line comments if present: // ...
	if r == '/' {
		r2, size2, err2 := lx.readRune()
		if err2 == nil && r2 == '/' {
			// consume until end of line, then get next token
			for {
				r3, _, err3 := lx.readRune()
				if err3 != nil {
					// EOF or other error ends comment; treat as EOF if EOF
					if err3 == io.EOF {
						return Token{Type: EOF, Pos: startPos}, nil
					}
					return Token{}, err3
				}
				if r3 == '\n' {
					break
				}
			}
			// After skipping comment, recurse to fetch next token
			return lx.NextToken()
		}
		if err2 == nil {
			lx.unreadRune(r2, size2)
		}
	}

	// Number literals (allow underscores between digits) and optional fractional part
	if unicode.IsDigit(r) {
		var sb strings.Builder
		sb.WriteRune(r)
		seenDot := false
		for {
			r2, size2, err2 := lx.readRune()
			if err2 != nil {
				if err2 == io.EOF {
					break
				}
				return Token{}, err2
			}
			if unicode.IsDigit(r2) || r2 == '_' {
				sb.WriteRune(r2)
				continue
			}
			if r2 == '.' && !seenDot {
				// ensure there's at least one digit after the dot
				r3, size3, err3 := lx.readRune()
				if err3 == nil && unicode.IsDigit(r3) {
					sb.WriteRune('.')
					sb.WriteRune(r3)
					seenDot = true
					continue
				}
				// not a fractional part; unread r3 (if any) and the '.' then stop number
				if err3 == nil {
					lx.unreadRune(r3, size3)
				}
				lx.unreadRune(r2, size2)
				break
			}
			// Not part of number
			if !unicode.IsSpace(r2) {
				lx.unreadRune(r2, size2)
			}
			break
		}
		lit := sb.String()
		return Token{Type: NUMBER, Lit: lit, Pos: startPos}, nil
	}

	// Strings
	if r == '"' {
		var sb strings.Builder
		for {
			r2, _, err2 := lx.readRune()
			if err2 != nil {
				if err2 == io.EOF {
					return Token{Type: ILLEGAL, Lit: "unterminated string", Pos: startPos}, nil
				}
				return Token{}, err2
			}
			if r2 == '"' {
				break
			}
			sb.WriteRune(r2)
		}
		return Token{Type: STRING, Lit: sb.String(), Pos: startPos}, nil
	}

	// Symbols and operators (including multi-char)
	switch r {
	case '=':
		// check '=='
		r2, size2, err2 := lx.readRune()
		if err2 == nil && r2 == '=' {
			return Token{Type: EQEQ, Lit: "==", Pos: startPos}, nil
		}
		if err2 == nil {
			lx.unreadRune(r2, size2)
		}
		return Token{Type: ASSIGN, Lit: "=", Pos: startPos}, nil
	case '!':
		// check '!='
		r2, size2, err2 := lx.readRune()
		if err2 == nil && r2 == '=' {
			return Token{Type: NEQ, Lit: "!=", Pos: startPos}, nil
		}
		if err2 == nil {
			lx.unreadRune(r2, size2)
		}
		return Token{Type: BANG, Lit: "!", Pos: startPos}, nil
	case '<':
		// check '<='
		r2, size2, err2 := lx.readRune()
		if err2 == nil && r2 == '=' {
			return Token{Type: LTE, Lit: "<=", Pos: startPos}, nil
		}
		if err2 == nil {
			lx.unreadRune(r2, size2)
		}
		return Token{Type: LANGLE, Lit: "<", Pos: startPos}, nil
	case '>':
		// check '>='
		r2, size2, err2 := lx.readRune()
		if err2 == nil && r2 == '=' {
			return Token{Type: GTE, Lit: ">=", Pos: startPos}, nil
		}
		if err2 == nil {
			lx.unreadRune(r2, size2)
		}
		return Token{Type: RANGLE, Lit: ">", Pos: startPos}, nil
	case '|':
		// check '||'
		r2, size2, err2 := lx.readRune()
		if err2 == nil && r2 == '|' {
			return Token{Type: OROR, Lit: "||", Pos: startPos}, nil
		}
		if err2 == nil {
			lx.unreadRune(r2, size2)
		}
		return Token{Type: PIPE, Lit: "|", Pos: startPos}, nil
	case '&':
		// '&&' only; single '&' is illegal in this grammar
		r2, size2, err2 := lx.readRune()
		if err2 == nil && r2 == '&' {
			return Token{Type: ANDAND, Lit: "&&", Pos: startPos}, nil
		}
		if err2 == nil {
			lx.unreadRune(r2, size2)
		}
		return Token{Type: ILLEGAL, Lit: "&", Pos: startPos}, nil
	case '(':
		return Token{Type: LPAREN, Lit: "(", Pos: startPos}, nil
	case ')':
		return Token{Type: RPAREN, Lit: ")", Pos: startPos}, nil
	case '{':
		return Token{Type: LBRACE, Lit: "{", Pos: startPos}, nil
	case '}':
		return Token{Type: RBRACE, Lit: "}", Pos: startPos}, nil
	case '[':
		return Token{Type: LBRACKET, Lit: "[", Pos: startPos}, nil
	case ']':
		return Token{Type: RBRACKET, Lit: "]", Pos: startPos}, nil
	case ',':
		return Token{Type: COMMA, Lit: ",", Pos: startPos}, nil
	case '.':
		return Token{Type: DOT, Lit: ".", Pos: startPos}, nil
	case ':':
		return Token{Type: COLON, Lit: ":", Pos: startPos}, nil
	case ';':
		return Token{Type: SEMICOLON, Lit: ";", Pos: startPos}, nil
	case '?':
		return Token{Type: QUESTION, Lit: "?", Pos: startPos}, nil
	case '+':
		return Token{Type: PLUS, Lit: "+", Pos: startPos}, nil
	case '-':
		return Token{Type: MINUS, Lit: "-", Pos: startPos}, nil
	case '*':
		return Token{Type: STAR, Lit: "*", Pos: startPos}, nil
	case '%':
		return Token{Type: PERCENT, Lit: "%", Pos: startPos}, nil
	case '/':
		return Token{Type: SLASH, Lit: "/", Pos: startPos}, nil
	}

	// Identifiers or keywords
	if isIdentStart(r) {
		var sb strings.Builder
		sb.WriteRune(r)
		for {
			r2, size2, err2 := lx.readRune()
			if err2 != nil {
				if err2 == io.EOF {
					break
				}
				return Token{}, err2
			}
			if !isIdentPart(r2) {
				if !unicode.IsSpace(r2) {
					lx.unreadRune(r2, size2)
				}
				break
			}
			sb.WriteRune(r2)
		}
		lit := sb.String()
		switch lit {
		case "package":
			return Token{Type: PACKAGE, Lit: lit, Pos: startPos}, nil
		case "import":
			return Token{Type: IMPORT, Lit: lit, Pos: startPos}, nil
		case "export":
			return Token{Type: EXPORT, Lit: lit, Pos: startPos}, nil
		case "def":
			return Token{Type: DEF, Lit: lit, Pos: startPos}, nil
		case "type":
			return Token{Type: TYPE, Lit: lit, Pos: startPos}, nil
		case "struct":
			return Token{Type: STRUCT, Lit: lit, Pos: startPos}, nil
		case "enum":
			return Token{Type: ENUM, Lit: lit, Pos: startPos}, nil
		case "impl":
			return Token{Type: IMPL, Lit: lit, Pos: startPos}, nil
		case "fun":
			return Token{Type: FUN, Lit: lit, Pos: startPos}, nil
		case "return":
			return Token{Type: RETURN, Lit: lit, Pos: startPos}, nil
		case "if":
			return Token{Type: IF, Lit: lit, Pos: startPos}, nil
		case "else":
			return Token{Type: ELSE, Lit: lit, Pos: startPos}, nil
		case "while":
			return Token{Type: WHILE, Lit: lit, Pos: startPos}, nil
		case "for":
			return Token{Type: FOR, Lit: lit, Pos: startPos}, nil
		case "in":
			return Token{Type: IN, Lit: lit, Pos: startPos}, nil
		case "handle":
			return Token{Type: HANDLE, Lit: lit, Pos: startPos}, nil
		case "Ok":
			return Token{Type: OK, Lit: lit, Pos: startPos}, nil
		case "Err":
			return Token{Type: ERR, Lit: lit, Pos: startPos}, nil
		case "as":
			return Token{Type: AS, Lit: lit, Pos: startPos}, nil
		case "true":
			return Token{Type: TRUE, Lit: lit, Pos: startPos}, nil
		case "false":
			return Token{Type: FALSE, Lit: lit, Pos: startPos}, nil
		}
		return Token{Type: IDENT, Lit: lit, Pos: startPos}, nil
	}

	// Unknown single rune => ILLEGAL
	return Token{
		Type: ILLEGAL,
		Lit:  string(r),
		Pos:  startPos,
	}, nil
}

func isIdentStart(r rune) bool {
	return r == '_' || unicode.IsLetter(r)
}

func isIdentPart(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}
