package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestEveryTokenType(t *testing.T) {
	line := 1
	tests := []struct {
		input string
		token Token
	}{
		{input: "(", token: newToken(LEFT_PAREN, "(", line)},
		{input: ")", token: newToken(RIGHT_PAREN, ")", line)},
		{input: "+", token: newToken(PLUS, "+", line)},
		{input: "-", token: newToken(MINUS, "-", line)},
		{input: "*", token: newToken(STAR, "*", line)},
		{input: "/", token: newToken(BACKSLASH, "/", line)},
		{input: ",", token: newToken(COMMMA, ",", line)},
		{input: "or", token: newToken(OR, "or", line)},
		{input: "and", token: newToken(AND, "and", line)},
		{input: "not", token: newToken(NOT, "not", line)},
		{input: "=", token: newToken(EQUAL, "=", line)},
		{input: ">", token: newToken(GREATER, ">", line)},
		{input: "<", token: newToken(LESS, "<", line)},
		{input: ">=", token: newToken(GREATER_EQUAL, ">=", line)},
		{input: "<=", token: newToken(LESS_EQUAL, "<=", line)},
		{input: "<>", token: newToken(LESS_GREATER, "<>", line)},
		{input: ":=", token: newToken(COLON_EQUAL, ":=", line)},
		{input: ";", token: newToken(SEMICOLON, ";", line)},
		{input: "x", token: newToken(IDENTIFIER, "x", line)},
		{input: "5", token: newToken(NUMBER, "5", line)},
		{input: "true", token: newToken(TRUE, "true", line)},
		{input: "false", token: newToken(FALSE, "false", line)},
		{input: "read", token: newToken(READ, "read", line)},
		{input: "write", token: newToken(WRITE, "write", line)},
		{input: "if", token: newToken(IF, "if", line)},
		{input: "else", token: newToken(ELSE, "else", line)},
		{input: "then", token: newToken(THEN, "then", line)},
		{input: "fi", token: newToken(FI, "fi", line)},
		{input: "while", token: newToken(WHILE, "while", line)},
		{input: "do", token: newToken(DO, "do", line)},
		{input: "end", token: newToken(END, "end", line)},
		{input: "", token: newToken(EOF, "", line)},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Single token test: %v", test.token), func(t *testing.T) {
			tokens, err := NewLexer(test.input).Scan()
			if err != nil {
				t.Fatalf("errors parsing tokens: %v", err)
			}
			assert.Equal(t, test.token, tokens[0])
		})
	}
}

func TestMultipleTokens(t *testing.T) {
	tests := []struct {
		file   string
		tokens []Token
	}{
		{
			"tests/gcd.txt",
			[]Token{
				newToken(READ, "read", 1),
				newToken(IDENTIFIER, "x", 1),
				newToken(COMMMA, ",", 1),
				newToken(IDENTIFIER, "y", 1),
				newToken(SEMICOLON, ";", 1),
				newToken(WHILE, "while", 2),
				newToken(IDENTIFIER, "x", 2),
				newToken(LESS_GREATER, "<>", 2),
				newToken(IDENTIFIER, "y", 2),
				newToken(DO, "do", 2),
				newToken(IF, "if", 3),
				newToken(IDENTIFIER, "x", 3),
				newToken(LESS_EQUAL, "<=", 3),
				newToken(IDENTIFIER, "y", 3),
				newToken(THEN, "then", 3),
				newToken(IDENTIFIER, "y", 4),
				newToken(COLON_EQUAL, ":=", 4),
				newToken(IDENTIFIER, "y", 4),
				newToken(MINUS, "-", 4),
				newToken(IDENTIFIER, "x", 4),
				newToken(SEMICOLON, ";", 4),
				newToken(ELSE, "else", 5),
				newToken(IDENTIFIER, "x", 6),
				newToken(COLON_EQUAL, ":=", 6),
				newToken(IDENTIFIER, "x", 6),
				newToken(MINUS, "-", 6),
				newToken(IDENTIFIER, "y", 6),
				newToken(SEMICOLON, ";", 6),
				newToken(FI, "fi", 7),
				newToken(SEMICOLON, ";", 7),
				newToken(END, "end", 8),
				newToken(SEMICOLON, ";", 8),
				newToken(WRITE, "write", 9),
				newToken(IDENTIFIER, "x", 9),
				newToken(COMMMA, ",", 9),
				newToken(IDENTIFIER, "y", 9),
				newToken(SEMICOLON, ";", 9),
				newToken(EOF, "\n", 10),
			},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Multiple token test for file: %v", test.file), func(t *testing.T) {
			b, err := os.ReadFile(test.file)
			if err != nil {
				t.Fatalf("error opening source file: %v", err)
			}
			tokens, errs := NewLexer(string(b)).Scan()
			if errs != nil {
				t.Fatalf("errors parsing tokens: %v", errs)
			}
			for i, token := range tokens {
				assert.Equal(t, test.tokens[i], token)
			}
		})
	}
}
