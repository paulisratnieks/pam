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
		t.Run(fmt.Sprintf("Test interpreter for program: %v", test.name), func(t *testing.T) {
			i := NewInterpreter(test.input)
			i.Interpret(test.ast)
			assert.Equal(t, test.output, i.output)
		})
	}
}
