package main

import (
	"fmt"
	"strings"
)

type Parser struct {
	source  []Token
	current int
	start   int
}

type Node interface {
	Print() string
}

type Expr interface {
	Node
	exprNode()
}

type IntExpr interface {
	Expr
	VisitableInt
	intExprNode()
}

type BoolExpr interface {
	Expr
	VisitableBool
	boolExprNode()
}

type Stmt interface {
	Node
	VisitableVoid
	stmtNode()
}

type RootNode struct {
	Stmts []Stmt
}

// Expressions
type BasicIntLit struct {
	Type  TokenType
	Value string
}

type BasicBoolLit struct {
	Type TokenType
}

type BinaryIntExpr struct {
	Left  IntExpr
	Op    TokenType
	Right IntExpr
}

type BinaryBoolExpr struct {
	Left  Expr
	Op    TokenType
	Right Expr
}

type UnaryBoolExpr struct {
	Op   TokenType
	Expr BoolExpr
}

type ParenIntExpr struct {
	Value IntExpr
}

type ParenBoolExpr struct {
	Value BoolExpr
}

// Statements
type ReadStmt struct {
	Vars []string
}

type WriteStmt struct {
	Vars []string
}

type AssignStmt struct {
	Id    string
	Value IntExpr
}

type CondStmt struct {
	Cond BoolExpr
	If   []Stmt
	Else []Stmt
}

type LoopStmt struct {
	Cond  BoolExpr
	Stmts []Stmt
}

func (e *RootNode) Accept(i *Interpreter) { i.visitRootNode(e) }
func (e *RootNode) Print() string {
	return fmt.Sprintf("%v", stmtsString(e.Stmts))
}
func (e *BasicIntLit) exprNode()                 {}
func (e *BasicIntLit) intExprNode()              {}
func (e *BasicIntLit) Accept(i *Interpreter) int { return i.visitBasicIntLit(e) }
func (e *BasicIntLit) Print() string {
	return fmt.Sprintf("%v:%v", e.Type, e.Value)
}
func (e *BasicBoolLit) exprNode()                  {}
func (e *BasicBoolLit) boolExprNode()              {}
func (e *BasicBoolLit) Accept(i *Interpreter) bool { return i.visitBasicBoolLit(e) }
func (e *BasicBoolLit) Print() string {
	return fmt.Sprintf("%v", e.Type)
}
func (e *BinaryIntExpr) exprNode()                 {}
func (e *BinaryIntExpr) intExprNode()              {}
func (e *BinaryIntExpr) Accept(i *Interpreter) int { return i.visitBinaryIntExpr(e) }
func (e *BinaryIntExpr) Print() string {
	return fmt.Sprintf("%v %v %v", e.Left.Print(), e.Op, e.Right.Print())
}
func (e *BinaryBoolExpr) exprNode()                  {}
func (e *BinaryBoolExpr) boolExprNode()              {}
func (e *BinaryBoolExpr) Accept(i *Interpreter) bool { return i.visitBinaryBoolExpr(e) }
func (e *BinaryBoolExpr) Print() string {
	return fmt.Sprintf("%v %v %v", e.Left.Print(), e.Op, e.Right.Print())
}
func (e *UnaryBoolExpr) exprNode()                  {}
func (e *UnaryBoolExpr) boolExprNode()              {}
func (e *UnaryBoolExpr) Accept(i *Interpreter) bool { return i.visitUnaryBoolExpr(e) }
func (e *UnaryBoolExpr) Print() string {
	return fmt.Sprintf("%v %v", e.Op, e.Expr.Print())
}
func (e *ParenIntExpr) exprNode()                 {}
func (e *ParenIntExpr) intExprNode()              {}
func (e *ParenIntExpr) Accept(i *Interpreter) int { return i.visitParenIntExpr(e) }
func (e *ParenIntExpr) Print() string {
	return fmt.Sprintf("( %v )", e.Value.Print())
}
func (e *ParenBoolExpr) exprNode()                  {}
func (e *ParenBoolExpr) boolExprNode()              {}
func (e *ParenBoolExpr) Accept(i *Interpreter) bool { return i.visitParenBoolExpr(e) }
func (e *ParenBoolExpr) Print() string {
	return fmt.Sprintf("( %v )", e.Value.Print())
}
func (e *AssignStmt) stmtNode()             {}
func (e *AssignStmt) Accept(i *Interpreter) { i.visitAssignStmt(e) }
func (e *AssignStmt) Print() string {
	return fmt.Sprintf("%v:=(%v)", e.Id, e.Value.Print())
}
func (e *ReadStmt) stmtNode()             {}
func (e *ReadStmt) Accept(i *Interpreter) { i.visitReadStmt(e) }
func (e *ReadStmt) Print() string {
	return fmt.Sprintf("read: %v", e.Vars)
}
func (e *WriteStmt) stmtNode()             {}
func (e *WriteStmt) Accept(i *Interpreter) { i.visitWriteStmt(e) }
func (e *WriteStmt) Print() string {
	return fmt.Sprintf("write: %v", e.Vars)
}
func (e *CondStmt) stmtNode()             {}
func (e *CondStmt) Accept(i *Interpreter) { i.visitCondStmt(e) }
func (e *CondStmt) Print() string {
	if e.Else != nil {
		return fmt.Sprintf("if (%v) then (%v) else (%v) fi", e.Cond.Print(), stmtsString(e.If), stmtsString(e.Else))
	}
	return fmt.Sprintf("if (%v) then (%v) fi", e.Cond.Print(), stmtsString(e.If))
}
func (e *LoopStmt) stmtNode()             {}
func (e *LoopStmt) Accept(i *Interpreter) { i.visitLoopStmt(e) }
func (e *LoopStmt) Print() string {
	return fmt.Sprintf("while (%v) do (%v) end", e.Cond.Print(), stmtsString(e.Stmts))
}

func stmtsString(stmts []Stmt) string {
	s := make([]string, 0)
	for _, stmt := range stmts {
		s = append(s, stmt.Print())
	}

	return strings.Join(s, "; ")
}

func newRootNode(stmts []Stmt) *RootNode {
	return &RootNode{stmts}
}
func newBinaryIntExpr(left IntExpr, op TokenType, right IntExpr) *BinaryIntExpr {
	return &BinaryIntExpr{left, op, right}
}
func newBinaryBoolExpr(left Expr, op TokenType, right Expr) *BinaryBoolExpr {
	return &BinaryBoolExpr{left, op, right}
}
func newUnaryBoolExpr(op TokenType, e BoolExpr) *UnaryBoolExpr {
	return &UnaryBoolExpr{op, e}
}
func newBasicIntLit(t TokenType, v string) *BasicIntLit {
	return &BasicIntLit{t, v}
}
func newBasicBoolLit(t TokenType) *BasicBoolLit {
	return &BasicBoolLit{t}
}
func newParenIntExpr(e IntExpr) *ParenIntExpr {
	return &ParenIntExpr{e}
}
func newParenBoolExpr(e BoolExpr) *ParenBoolExpr {
	return &ParenBoolExpr{e}
}
func newAssignStmt(id string, e IntExpr) *AssignStmt {
	return &AssignStmt{id, e}
}
func newReadStmt(v []string) *ReadStmt {
	return &ReadStmt{v}
}
func newWriteStmt(v []string) *WriteStmt {
	return &WriteStmt{v}
}
func newCondStmt(cond BoolExpr, s ...[]Stmt) *CondStmt {
	if len(s) == 2 {
		return &CondStmt{cond, s[0], s[1]}
	}

	return &CondStmt{Cond: cond, If: s[0]}
}
func newLoopStmt(cond BoolExpr, s []Stmt) *LoopStmt {
	return &LoopStmt{cond, s}
}

func NewParser(s []Token) *Parser {
	return &Parser{source: s}
}

func (p *Parser) Parse() *RootNode {
	return newRootNode(p.series())
}

func (p *Parser) peek() Token {
	return p.source[p.current]
}

func (p *Parser) match(t TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	if p.peek().Type != t {
		return false
	}
	p.current++
	return true
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == Eof
}

func (p *Parser) nextToken() Token {
	if p.isAtEnd() {
		panic("unexpected end of file")
	}
	t := p.source[p.current]
	p.current++
	return t
}

func (p *Parser) consume(t TokenType) {
	if p.peek().Type != t {
		panic(fmt.Sprintf("expecting token of type: %v on token: %v but got: %v", t, p.current+1, p.peek().Type))
	}
	p.current++
}

func (p *Parser) previous() Token {
	return p.source[p.current-1]
}

func (p *Parser) progr() []Stmt {
	return p.series()
}

func (p *Parser) series() []Stmt {
	s := make([]Stmt, 0)
	s = append(s, p.statement())
	p.consume(Semicolon)

	t := p.peek().Type
	for !p.isAtEnd() && !(t == End || t == Else || t == Fi) {
		s = append(s, p.statement())
		p.consume(Semicolon)
		t = p.peek().Type
	}

	return s
}

func (p *Parser) statement() Stmt {
	switch p.peek().Type {
	case Identifier:
		return p.assignStmt()
	case Read:
		return p.readStmt()
	case Write:
		return p.writeStmt()
	case If:
		return p.condStmt()
	default:
		return p.loopStmt()
	}
}

func (p *Parser) assignStmt() Stmt {
	id := p.nextToken()
	p.consume(ColonEqual)

	return newAssignStmt(id.Lexeme, p.expression())
}

func (p *Parser) readStmt() Stmt {
	p.consume(Read)

	return newReadStmt(p.varList())
}

func (p *Parser) writeStmt() Stmt {
	p.consume(Write)

	return newWriteStmt(p.varList())
}

func (p *Parser) condStmt() Stmt {
	p.consume(If)
	cond := p.logical()
	p.consume(Then)
	left := p.series()

	var right []Stmt
	if p.match(Else) {
		right = p.series()
	}
	p.consume(Fi)

	if right != nil {
		return newCondStmt(cond, left, right)
	}

	return newCondStmt(cond, left)
}

func (p *Parser) loopStmt() Stmt {
	p.consume(While)
	cond := p.logical()
	p.consume(Do)
	series := p.series()
	p.consume(End)

	return newLoopStmt(cond, series)
}

func (p *Parser) varList() []string {
	vars := make([]string, 0)
	t := p.nextToken()
	vars = append(vars, t.Lexeme)

	for p.match(Comma) {
		vars = append(vars, p.nextToken().Lexeme)
	}

	return vars
}

func (p *Parser) logical() BoolExpr {
	return p.logicalOr()
}

func (p *Parser) logicalOr() BoolExpr {
	e := p.logicalAnd()

	for p.match(Or) {
		op := p.previous().Type
		right := p.logicalAnd()
		e = newBinaryBoolExpr(e, op, right)
	}

	return e
}

func (p *Parser) logicalAnd() BoolExpr {
	e := p.logicalNot()

	for p.match(And) {
		op := p.previous().Type
		right := p.logicalNot()
		e = newBinaryBoolExpr(e, op, right)
	}

	return e
}

func (p *Parser) logicalNot() BoolExpr {
	if p.match(Not) {
		return newUnaryBoolExpr(p.previous().Type, p.logical())
	}

	return p.logicalRelation()
}

func (p *Parser) logicalRelation() BoolExpr {
	t := p.peek()
	if t.Type == True || t.Type == False || t.Type == LeftParen {
		return p.logicalElem()
	} else {
		return newBinaryBoolExpr(p.expression(), p.nextToken().Type, p.expression())
	}
}

func (p *Parser) logicalElem() BoolExpr {
	t := p.nextToken()

	if t.Type == True || t.Type == False {
		return newBasicBoolLit(t.Type)
	} else {
		e := p.logical()
		p.consume(RightParen)
		return newParenBoolExpr(e)
	}
}

func (p *Parser) expression() IntExpr {
	exp := p.term()

	for p.match(Plus) || p.match(Minus) {
		op := p.previous().Type
		right := p.term()
		exp = newBinaryIntExpr(exp, op, right)
	}

	return exp
}

func (p *Parser) term() IntExpr {
	elem := p.element()

	for p.match(Star) || p.match(Backslash) {
		op := p.previous().Type
		right := p.element()
		elem = newBinaryIntExpr(elem, op, right)
	}

	return elem
}

func (p *Parser) element() IntExpr {
	t := p.nextToken()
	switch t.Type {
	case Number, Identifier:
		return newBasicIntLit(t.Type, t.Lexeme)
	default:
		elem := newParenIntExpr(p.expression())
		p.consume(RightParen)
		return elem
	}
}
