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
			newRootNode([]Stmt{
				newReadStmt([]string{"x", "y"}),
				newLoopStmt(
					newBinaryBoolExpr(
						newBasicIntLit(Identifier, "x"),
						LessGreater,
						newBasicIntLit(Identifier, "y"),
					),
					[]Stmt{
						newCondStmt(
							newBinaryBoolExpr(
								newBasicIntLit(Identifier, "x"),
								LessEqual,
								newBasicIntLit(Identifier, "y"),
							),
							[]Stmt{
								newAssignStmt(
									"y",
									newBinaryIntExpr(
										newBasicIntLit(Identifier, "y"),
										Minus,
										newBasicIntLit(Identifier, "x"),
									),
								),
							},
							[]Stmt{
								newAssignStmt(
									"x",
									newBinaryIntExpr(
										newBasicIntLit(Identifier, "x"),
										Minus,
										newBasicIntLit(Identifier, "y"),
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
