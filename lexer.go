package main

import "fmt"

type TokenType string

const (
	LeftParen    = "LeftParen"
	RightParen   = "RightParen"
	Plus         = "Plus"
	Minus        = "Minus"
	Star         = "Star"
	Backslash    = "Backslash"
	Comma        = "Comma"
	Or           = "Or"
	And          = "And"
	Not          = "Not"
	Equal        = "Equal"
	Greater      = "Greater"
	Less         = "Less"
	GreaterEqual = "GreaterEqual"
	LessEqual    = "LessEqual"
	LessGreater  = "LessGreater"
	ColonEqual   = "ColonEqual"
	Semicolon    = "Semicolon"
	Identifier   = "Identifier"
	Number       = "Number"
	True         = "True"
	False        = "False"
	Read         = "Read"
	Write        = "Write"
	If           = "If"
	Else         = "Else"
	Then         = "Then"
	Fi           = "Fi"
	While        = "While"
	Do           = "Do"
	End          = "End"
	Eof          = "Eof"
)

var keywords = map[string]TokenType{
	"if":    If,
	"else":  Else,
	"then":  Then,
	"fi":    Fi,
	"while": While,
	"do":    Do,
	"end":   End,
	"write": Write,
	"read":  Read,
	"or":    Or,
	"and":   And,
	"not":   Not,
	"true":  True,
	"false": False,
}

type Token struct {
	Type   TokenType
	Lexeme string
	Line   int
}

type Lexer struct {
	source  string
	current int
	start   int
	line    int
	tokens  []Token
	errors  []error
}

func isDigit(c uint8) bool {
	return c >= '0' && c <= '9'
}

func isAlphaNum(c uint8) bool {
	return c >= 'a' && c <= 'z' ||
		c >= 'A' && c <= 'Z' ||
		c >= '0' && c <= '9'
}

func (l *Lexer) Scan() ([]Token, []error) {
	for !l.isAtEnd() {
		l.start = l.current

		switch char := l.nextChar(); char {
		case '(':
			l.addToken(LeftParen)
		case ')':
			l.addToken(RightParen)
		case '+':
			l.addToken(Plus)
		case '-':
			l.addToken(Minus)
		case '*':
			l.addToken(Star)
		case '/':
			l.addToken(Backslash)
		case ',':
			l.addToken(Comma)
		case ';':
			l.addToken(Semicolon)
		case ':':
			l.nextChar()
			l.addToken(ColonEqual)
		case '>':
			if l.match('=') {
				l.addToken(GreaterEqual)
			} else {
				l.addToken(Greater)
			}
		case '<':
			if l.match('=') {
				l.addToken(LessEqual)
			} else if l.match('>') {
				l.addToken(LessGreater)
			} else {
				l.addToken(Less)
			}
		case '=':
			l.addToken(Equal)
		case '\n':
			l.line++
		case ' ':
		case '\t':
		case '\r':
			break
		default:
			if isDigit(char) {
				for isDigit(l.peek()) {
					l.nextChar()
				}
				l.addToken(Number)
			} else if isAlphaNum(char) {
				for isAlphaNum(l.peek()) {
					l.nextChar()
				}
				kw := keywords[l.source[l.start:l.current]]
				if kw != "" {
					l.addToken(kw)
				} else {
					l.addToken(Identifier)
				}
			} else {
				l.addError(fmt.Errorf("unexpected character '%c' at line %d column ", char, l.line))
			}
		}
	}

	l.addToken(Eof)

	return l.tokens, l.errors
}

func (l *Lexer) addToken(t TokenType) {
	l.tokens = append(l.tokens, newToken(t, l.source[l.start:l.current], l.line))
}

func newToken(t TokenType, le string, li int) Token {
	return Token{t, le, li}
}

func (l *Lexer) addError(e error) {
	l.errors = append(l.errors, e)
}

func (l *Lexer) isAtEnd() bool {
	return l.current >= len(l.source)
}

func (l *Lexer) match(c uint8) bool {
	if l.isAtEnd() {
		return false
	}
	if l.source[l.current] != c {
		return false
	}
	l.current++
	return true
}

func (l *Lexer) peek() uint8 {
	if l.isAtEnd() {
		return 0
	}
	return l.source[l.current]
}

func (l *Lexer) nextChar() uint8 {
	char := l.source[l.current]
	l.current++

	return char
}

func NewLexer(s string) *Lexer {
	return &Lexer{source: s, line: 1}
}
