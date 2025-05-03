package main

import (
	"fmt"
	"strconv"
)

const (
	integer = "int"
	boolean = "bool"
)

type VisitableVoid interface {
	Accept(i *Interpreter)
}

type VisitableInt interface {
	Accept(i *Interpreter) int
}

type VisitableBool interface {
	Accept(i *Interpreter) bool
}

type Interpreter struct {
	memory  map[string]int
	input   []int
	current int
	output  []int
}

func (i *Interpreter) nextInput() int {
	if i.current >= len(i.input) {
		i.error(fmt.Sprintf("not enough input tokens at nr: %v", i.current+1))
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
	s.Accept(i)
}

func (i *Interpreter) visitRootNode(n *RootNode) {
	for _, stmt := range n.Stmts {
		stmt.Accept(i)
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
		v, ok := i.memory[k]
		if !ok {
			i.error(fmt.Sprintf("undefined variable: %v", k))
		}
		i.output = append(i.output, v)
	}
}

func (i *Interpreter) visitAssignStmt(s *AssignStmt) {
	i.memory[s.Id] = s.Value.Accept(i)
}

func (i *Interpreter) visitCondStmt(s *CondStmt) {
	if s.Cond.Accept(i) {
		for _, stmt := range s.If {
			stmt.Accept(i)
		}
	} else {
		for _, stmt := range s.Else {
			stmt.Accept(i)
		}
	}
}

func (i *Interpreter) visitLoopStmt(s *LoopStmt) {
	cond := s.Cond.Accept(i)
	for cond {
		for _, stmt := range s.Stmts {
			stmt.Accept(i)
		}
		cond = s.Cond.Accept(i)
	}
}

func (i *Interpreter) visitBasicIntLit(e *BasicIntLit) int {
	if e.Type == Identifier {
		return i.memory[e.Value]
	}

	v, err := strconv.Atoi(e.Value)
	if err != nil {
		i.error(err.Error())
	}
	return v
}

func (i *Interpreter) visitBasicBoolLit(e *BasicBoolLit) bool {
	return e.Type == True
}

func (i *Interpreter) visitBinaryIntExpr(e *BinaryIntExpr) int {
	left := e.Left.Accept(i)
	right := e.Right.Accept(i)

	switch e.Op {
	case Plus:
		return left + right
	case Minus:
		return left - right
	case Star:
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
		left := leftBool.Accept(i)
		right := rightBool.Accept(i)

		switch e.Op {
		case And:
			return left && right
		default:
			return left || right
		}
	} else if isLeftInt && isRightInt {
		left := leftInt.Accept(i)
		right := rightInt.Accept(i)

		switch e.Op {
		case Equal:
			return left == right
		case Less:
			return left < right
		case Greater:
			return left > right
		case LessEqual:
			return left <= right
		case GreaterEqual:
			return left >= right
		default:
			return left != right
		}
	} else {
		leftType := integer
		if isLeftBool {
			leftType = boolean
		}
		rightType := integer
		if isRightBool {
			rightType = boolean
		}
		i.error(fmt.Sprintf("type mismatch between expressions: %v %v %v", leftType, e.Op, rightType))
		return false
	}
}

func (i *Interpreter) visitUnaryBoolExpr(e *UnaryBoolExpr) bool {
	return !e.Expr.Accept(i)
}

func (i *Interpreter) visitParenIntExpr(e *ParenIntExpr) int {
	return e.Value.Accept(i)
}

func (i *Interpreter) visitParenBoolExpr(e *ParenBoolExpr) bool {
	return e.Value.Accept(i)
}
