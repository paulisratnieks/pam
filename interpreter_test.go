package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInterpreter(t *testing.T) {
	tests := []struct {
		name   string
		input  []int
		output []int
		ast    *RootNode
	}{
		{
			"gcd",
			[]int{14, 35},
			[]int{7, 7},
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
		t.Run(fmt.Sprintf("Test interpreter for program: %v", test.name), func(t *testing.T) {
			i := NewInterpreter(test.input)
			i.Interpret(test.ast)
			assert.Equal(t, test.output, i.output)
		})
	}
}
