package main

import "fmt"

type TokenType string

const (
	LEFT_PAREN    = "LEFT_PAREN"
	RIGHT_PAREN   = "RIGHT_PAREN"
	PLUS          = "PLUS"
	MINUS         = "MINUS"
	STAR          = "STAR"
	BACKSLASH     = "BACKSLASH"
	COMMMA        = "COMMMA"
	OR            = "OR"
	AND           = "AND"
	NOT           = "NOT"
	EQUAL         = "EQUAL"
	GREATER       = "GREATER"
	LESS          = "LESS"
	GREATER_EQUAL = "GREATER_EQUAL"
	LESS_EQUAL    = "LESS_EQUAL"
	LESS_GREATER  = "LESS_GREATER"
	COLON_EQUAL   = "COLON_EQUAL"
	SEMICOLON     = "SEMICOLON"
	IDENTIFIER    = "IDENTIFIER"
	NUMBER        = "NUMBER"
	TRUE          = "TRUE"
	FALSE         = "FALSE"
	READ          = "READ"
	WRITE         = "WRITE"
	IF            = "IF"
	ELSE          = "ELSE"
	THEN          = "THEN"
	FI            = "FI"
	WHILE         = "WHILE"
	DO            = "DO"
	END           = "END"
	EOF           = "EOF"
)

var keywords = map[string]TokenType{
	"if":    IF,
	"else":  ELSE,
	"then":  THEN,
	"fi":    FI,
	"while": WHILE,
	"do":    DO,
	"end":   END,
	"write": WRITE,
	"read":  READ,
	"or":    OR,
	"and":   AND,
	"not":   NOT,
	"true":  TRUE,
	"false": FALSE,
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
			l.addToken(LEFT_PAREN)
		case ')':
			l.addToken(RIGHT_PAREN)
		case '+':
			l.addToken(PLUS)
		case '-':
			l.addToken(MINUS)
		case '*':
			l.addToken(STAR)
		case '/':
			l.addToken(BACKSLASH)
		case ',':
			l.addToken(COMMMA)
		case ';':
			l.addToken(SEMICOLON)
		case ':':
			l.nextChar()
			l.addToken(COLON_EQUAL)
		case '>':
			if l.match('=') {
				l.addToken(GREATER_EQUAL)
			} else {
				l.addToken(GREATER)
			}
		case '<':
			if l.match('=') {
				l.addToken(LESS_EQUAL)
			} else if l.match('>') {
				l.addToken(LESS_GREATER)
			} else {
				l.addToken(LESS)
			}
		case '=':
			l.addToken(EQUAL)
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
				l.addToken(NUMBER)
			} else if isAlphaNum(char) {
				for isAlphaNum(l.peek()) {
					l.nextChar()
				}
				kw := keywords[l.source[l.start:l.current]]
				if kw != "" {
					l.addToken(kw)
				} else {
					l.addToken(IDENTIFIER)
				}
			} else {
				l.addError(fmt.Errorf("unexpected character '%c' at line %d column ", char, l.line))
			}
		}
	}

	l.addToken(EOF)

	return l.tokens, l.errors
}

func (l *Lexer) addToken(t TokenType) {
	l.tokens = append(l.tokens, Token{Type: t, Lexeme: l.source[l.start:l.current], Line: l.line})
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
