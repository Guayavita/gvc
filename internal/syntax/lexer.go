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

	// Identifiers
	IDENT

	// Literals
	NUMBER

	// Symbols
	ASSIGN // '='
	PIPE   // '|'

	// Other Keywords
	DEF
)

func (t TokenType) String() string {
	switch t {
	case EOF:
		return "EOF"
	case ILLEGAL:
		return "ILLEGAL"
	case PACKAGE:
		return "PACKAGE"
	case IDENT:
		return "IDENT"
	case NUMBER:
		return "NUMBER"
	case ASSIGN:
		return "ASSIGN"
	case PIPE:
		return "PIPE"
	case DEF:
		return "DEF"
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
		// Skip whitespace
		if unicode.IsSpace(r) {
			continue
		}
		// Line comment: // ...
		if r == '/' {
			r2, size2, err2 := lx.readRune()
			if err2 == nil && r2 == '/' {
				// consume until end of line
				for {
					r3, _, err3 := lx.readRune()
					if err3 != nil {
						if err3 == io.EOF {
							return nil
						}
						return err3
					}
					if r3 == '\n' {
						break
					}
				}
				continue
			}
			// Not a comment; unread second rune and treat '/' as ILLEGAL here
			if err2 == nil {
				lx.unreadRune(r2, size2)
			}
			lx.unreadRune(r, size)
			return nil
		}
		// Non-whitespace, non-comment start: unread it and return
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
			lx.unreadRune(r2, size2)
			break
		}
		lit := sb.String()
		return Token{Type: NUMBER, Lit: lit, Pos: startPos}, nil
	}

	// Single-char symbols
	switch r {
	case '=':
		return Token{Type: ASSIGN, Lit: "=", Pos: startPos}, nil
	case '|':
		return Token{Type: PIPE, Lit: "|", Pos: startPos}, nil
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
				lx.unreadRune(r2, size2)
				break
			}
			sb.WriteRune(r2)
		}
		lit := sb.String()
		if lit == "package" {
			return Token{Type: PACKAGE, Lit: lit, Pos: startPos}, nil
		}
		if lit == "def" {
			return Token{Type: DEF, Lit: lit, Pos: startPos}, nil
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
