package main

import (
	"fmt"
	"strconv"
)

const (
	Int  = "int"
	Bool = "bool"
)

type VisitableVoid interface {
	accept(i *Interpreter)
}

type VisitableInt interface {
	accept(i *Interpreter) int
}

type VisitableBool interface {
	accept(i *Interpreter) bool
}

type Interpreter struct {
	memory  map[string]int
	input   []int
	current int
	output  []int
}

func (i *Interpreter) nextInput() int {
	if i.current >= len(i.input) {
		i.error(fmt.Sprintf("enough input tokens at nr: %v", i.current+1))
	}
	input := i.input[i.current]
	i.current++
	return input
}

func (i *Interpreter) error(msg string) {
	panic(msg)
}

func NewInterpreter(input []int) *Interpreter {
	return &Interpreter{
		memory: make(map[string]int),
		input:  input,
	}
}

func (i *Interpreter) Interpret(s *RootNode) {
	s.accept(i)
}

func (i *Interpreter) visitRootNode(n *RootNode) {
	for _, stmt := range n.Stmts {
		stmt.accept(i)
	}
}

func (i *Interpreter) visitReadStmt(s *ReadStmt) {
	for _, k := range s.Vars {
		v := i.nextInput()
		i.memory[k] = v
	}
}

func (i *Interpreter) visitWriteStmt(s *WriteStmt) {
	for _, k := range s.Vars {
		v := i.memory[k]
		if v == 0 {
			i.error(fmt.Sprintf("undefined variable: %v", k))
		}
		i.output = append(i.output, v)
	}
}

func (i *Interpreter) visitAssignStmt(s *AssignStmt) {
	i.memory[s.Id] = s.Value.accept(i)
}

func (i *Interpreter) visitCondStmt(s *CondStmt) {
	if s.Cond.accept(i) {
		for _, stmt := range s.If {
			stmt.accept(i)
		}
	} else {
		for _, stmt := range s.Else {
			stmt.accept(i)
		}
	}
}

func (i *Interpreter) visitLoopStmt(s *LoopStmt) {
	cond := s.Cond.accept(i)
	for cond {
		for _, stmt := range s.Stmts {
			stmt.accept(i)
		}
		cond = s.Cond.accept(i)
	}
}

func (i *Interpreter) visitBasicIntLit(e *BasicIntLit) int {
	if e.Type == IDENTIFIER {
		return i.memory[e.Value]
	}

	v, err := strconv.Atoi(e.Value)
	if err != nil {
		i.error(err.Error())
	}
	return v
}

func (i *Interpreter) visitBasicBoolLit(e *BasicBoolLit) bool {
	return e.Type == TRUE
}

func (i *Interpreter) visitBinaryIntExpr(e *BinaryIntExpr) int {
	left := e.Left.accept(i)
	right := e.Right.accept(i)

	switch e.Op {
	case PLUS:
		return left + right
	case MINUS:
		return left - right
	case STAR:
		return left * right
	default:
		return left / right
	}
}

func (i *Interpreter) visitBinaryBoolExpr(e *BinaryBoolExpr) bool {
	leftBool, isLeftBool := e.Left.(BoolExpr)
	leftInt, isLeftInt := e.Left.(IntExpr)
	if !isLeftBool && !isLeftInt {
		i.error(fmt.Sprintf("unknown left expression type"))
	}

	rightBool, isRightBool := e.Right.(BoolExpr)
	rightInt, isRightInt := e.Right.(IntExpr)
	if !isLeftBool && !isLeftInt {
		i.error(fmt.Sprintf("unknown left expression type"))
	}

	if isLeftBool && isRightBool {
		left := leftBool.accept(i)
		right := rightBool.accept(i)

		switch e.Op {
		case AND:
			return left && right
		default:
			return left || right
		}
	} else if isLeftInt && isRightInt {
		left := leftInt.accept(i)
		right := rightInt.accept(i)

		switch e.Op {
		case EQUAL:
			return left == right
		case LESS:
			return left < right
		case GREATER:
			return left > right
		case LESS_EQUAL:
			return left <= right
		case GREATER_EQUAL:
			return left >= right
		default:
			return left != right
		}
	} else {
		leftType := Int
		if isLeftBool {
			leftType = Bool
		}
		rightType := Int
		if isRightBool {
			rightType = Bool
		}
		i.error(fmt.Sprintf("type mismatch between expressions: %v %v %v", leftType, e.Op, rightType))
		return false
	}
}

func (i *Interpreter) visitUnaryBoolExpr(e *UnaryBoolExpr) bool {
	return !e.Expr.accept(i)
}

func (i *Interpreter) visitParenIntExpr(e *ParenIntExpr) int {
	return e.Value.accept(i)
}

func (i *Interpreter) visitParenBoolExpr(e *ParenBoolExpr) bool {
	return e.Value.accept(i)
}
