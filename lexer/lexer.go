package lexer

import (
	"avm/token"
	"strings"
)

type Lexer struct {
	in      string
	Pos     int // current position (points to current char)
	readPos int // current reading position in input. Always point to the next char in the input
	ch      byte
}

func New(in string) *Lexer {
	l := &Lexer{in: in}
	l.scan()
	return l
}

// scan find the next char
func (l *Lexer) scan() {
	if l.readPos >= len(l.in) {
		l.ch = 0
	} else {
		l.ch = l.in[l.readPos]
	}

	l.Pos = l.readPos
	l.readPos++
}

// NextToken returns the token
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.scanIgnoreWhiteSpace()

	switch l.ch {
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ';':
		if l.peekChar() == ';' {
			ch := l.ch
			l.scan()
			lit := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EOI, Literal: lit}
		} else {
			tok = newToken(token.SEMICOLON, l.ch)
		}
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '/':
		tok = newToken(token.SLASH, l.ch)
	default:
		if isLetter(l.ch) {
			tok.Literal = l.ScanIdent()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Literal = l.scanNumber()
			if strings.Contains(tok.Literal, ".") {
				tok.Type = token.FLOAT_NUM
			} else {
				tok.Type = token.INT
			}
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}

	}

	l.scan()
	return tok

}

func (l *Lexer) scanNumber() string {
	pos := l.Pos
	for isDigit(l.ch) {
		l.scan()
	}

	return l.in[pos:l.Pos]
}

func (l *Lexer) peekChar() byte {
	if l.readPos >= len(l.in) {
		return 0
	}

	return l.in[l.readPos]
}

// ScanIdent read until it encounters non-letter character
func (l *Lexer) ScanIdent() string {
	pos := l.Pos
	for isLetter(l.ch) {
		l.scan()
	}

	for isDigit(l.ch) {
		l.scan()
	}

	return l.in[pos:l.Pos]
}

// isLetter check if the given parameter is a letter.
func isLetter(ch byte) bool {
	return ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z'
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9' || ch == '.'
}

// newToken return a new Token
func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{tokenType, string(ch)}
}

func (l *Lexer) scanIgnoreWhiteSpace() {
	for l.ch == ' ' || l.ch == '\r' || l.ch == '\t' || l.ch == '\n' {
		l.scan()
	}

}
