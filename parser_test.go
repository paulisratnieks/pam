package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParser(t *testing.T) {
	tests := []struct {
		name   string
		tokens []Token
		ast    *RootNode
	}{
		{
			"gcd",
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
			newRootNode([]Stmt{
				newReadStmt([]string{"x", "y"}),
				newLoopStmt(
					newBinaryBoolExpr(
						newBasicIntLit(IDENTIFIER, "x"),
						LESS_GREATER,
						newBasicIntLit(IDENTIFIER, "y"),
					),
					[]Stmt{
						newCondStmt(
							newBinaryBoolExpr(
								newBasicIntLit(IDENTIFIER, "x"),
								LESS_EQUAL,
								newBasicIntLit(IDENTIFIER, "y"),
							),
							[]Stmt{
								newAssignStmt(
									"y",
									newBinaryIntExpr(
										newBasicIntLit(IDENTIFIER, "y"),
										MINUS,
										newBasicIntLit(IDENTIFIER, "x"),
									),
								),
							},
							[]Stmt{
								newAssignStmt(
									"x",
									newBinaryIntExpr(
										newBasicIntLit(IDENTIFIER, "x"),
										MINUS,
										newBasicIntLit(IDENTIFIER, "y"),
									),
								),
							},
						),
					},
				),
				newWriteStmt([]string{"x", "y"}),
			}),
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Test parser for program: %v", test.name), func(t *testing.T) {
			ast := NewParser(test.tokens).Parse()
			assert.Equal(t, test.ast, ast)
		})
	}
}
