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
		{input: "(", token: newToken(LeftParen, "(", line)},
		{input: ")", token: newToken(RightParen, ")", line)},
		{input: "+", token: newToken(Plus, "+", line)},
		{input: "-", token: newToken(Minus, "-", line)},
		{input: "*", token: newToken(Star, "*", line)},
		{input: "/", token: newToken(Backslash, "/", line)},
		{input: ",", token: newToken(Comma, ",", line)},
		{input: "or", token: newToken(Or, "or", line)},
		{input: "and", token: newToken(And, "and", line)},
		{input: "not", token: newToken(Not, "not", line)},
		{input: "=", token: newToken(Equal, "=", line)},
		{input: ">", token: newToken(Greater, ">", line)},
		{input: "<", token: newToken(Less, "<", line)},
		{input: ">=", token: newToken(GreaterEqual, ">=", line)},
		{input: "<=", token: newToken(LessEqual, "<=", line)},
		{input: "<>", token: newToken(LessGreater, "<>", line)},
		{input: ":=", token: newToken(ColonEqual, ":=", line)},
		{input: ";", token: newToken(Semicolon, ";", line)},
		{input: "x", token: newToken(Identifier, "x", line)},
		{input: "5", token: newToken(Number, "5", line)},
		{input: "true", token: newToken(True, "true", line)},
		{input: "false", token: newToken(False, "false", line)},
		{input: "read", token: newToken(Read, "read", line)},
		{input: "write", token: newToken(Write, "write", line)},
		{input: "if", token: newToken(If, "if", line)},
		{input: "else", token: newToken(Else, "else", line)},
		{input: "then", token: newToken(Then, "then", line)},
		{input: "fi", token: newToken(Fi, "fi", line)},
		{input: "while", token: newToken(While, "while", line)},
		{input: "do", token: newToken(Do, "do", line)},
		{input: "end", token: newToken(End, "end", line)},
		{input: "", token: newToken(Eof, "", line)},
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
				newToken(Read, "read", 1),
				newToken(Identifier, "x", 1),
				newToken(Comma, ",", 1),
				newToken(Identifier, "y", 1),
				newToken(Semicolon, ";", 1),
				newToken(While, "while", 2),
				newToken(Identifier, "x", 2),
				newToken(LessGreater, "<>", 2),
				newToken(Identifier, "y", 2),
				newToken(Do, "do", 2),
				newToken(If, "if", 3),
				newToken(Identifier, "x", 3),
				newToken(LessEqual, "<=", 3),
				newToken(Identifier, "y", 3),
				newToken(Then, "then", 3),
				newToken(Identifier, "y", 4),
				newToken(ColonEqual, ":=", 4),
				newToken(Identifier, "y", 4),
				newToken(Minus, "-", 4),
				newToken(Identifier, "x", 4),
				newToken(Semicolon, ";", 4),
				newToken(Else, "else", 5),
				newToken(Identifier, "x", 6),
				newToken(ColonEqual, ":=", 6),
				newToken(Identifier, "x", 6),
				newToken(Minus, "-", 6),
				newToken(Identifier, "y", 6),
				newToken(Semicolon, ";", 6),
				newToken(Fi, "fi", 7),
				newToken(Semicolon, ";", 7),
				newToken(End, "end", 8),
				newToken(Semicolon, ";", 8),
				newToken(Write, "write", 9),
				newToken(Identifier, "x", 9),
				newToken(Comma, ",", 9),
				newToken(Identifier, "y", 9),
				newToken(Semicolon, ";", 9),
				newToken(Eof, "\n", 10),
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
