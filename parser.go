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
	print() string
}

type Expr interface {
	Node
	exprNode()
}

type Stmt interface {
	Node
	stmtNode()
}

type RootNode struct {
	Stmts []Stmt
}

// Expressions
type BasicLit struct {
	Type  TokenType
	Value string
}

type BinaryExpr struct {
	Left  Expr
	Op    TokenType
	Right Expr
}

type UnaryExpr struct {
	Op   TokenType
	Expr Expr
}

type ParenExpr struct {
	Value Expr
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
	Value Expr
}

type CondStmt struct {
	Cond Expr
	If   []Stmt
	Else []Stmt
}

type LoopStmt struct {
	Cond  Expr
	Stmts []Stmt
}

func (e *RootNode) print() string {
	return fmt.Sprintf("%v", stmtsString(e.Stmts))
}
func (e *BasicLit) exprNode() {}
func (e *BasicLit) print() string {
	return fmt.Sprintf("%v:%v", e.Type, e.Value)
}
func (e *BinaryExpr) exprNode() {}
func (e *BinaryExpr) print() string {
	return fmt.Sprintf("%v %v %v", e.Left.print(), e.Op, e.Right.print())
}
func (e *UnaryExpr) exprNode() {}
func (e *UnaryExpr) print() string {
	return fmt.Sprintf("%v %v", e.Op, e.Expr.print())
}
func (e *ParenExpr) exprNode() {}
func (e *ParenExpr) print() string {
	return fmt.Sprintf("( %v )", e.Value.print())
}
func (e *AssignStmt) stmtNode() {}
func (e *AssignStmt) print() string {
	return fmt.Sprintf("%v:=%v", e.Id, e.Value.print())
}
func (e *ReadStmt) stmtNode() {}
func (e *ReadStmt) print() string {
	return fmt.Sprintf("read: %v", e.Vars)
}
func (e *WriteStmt) stmtNode() {}
func (e *WriteStmt) print() string {
	return fmt.Sprintf("write: %v", e.Vars)
}
func (e *CondStmt) stmtNode() {}
func (e *CondStmt) print() string {
	if e.Else != nil {
		return fmt.Sprintf("if %v then %v else %v fi", e.Cond.print(), stmtsString(e.If), stmtsString(e.Else))
	}
	return fmt.Sprintf("if %v then %v fi", e.Cond.print(), stmtsString(e.If))
}
func (e *LoopStmt) stmtNode() {}
func (e *LoopStmt) print() string {
	return fmt.Sprintf("while %v do %v end", e.Cond.print(), stmtsString(e.Stmts))
}

func stmtsString(stmts []Stmt) string {
	s := make([]string, 0)
	for _, stmt := range stmts {
		s = append(s, stmt.print())
	}

	return strings.Join(s, "; ")
}

func NewRootNode(stmts []Stmt) *RootNode {
	return &RootNode{stmts}
}

func NewBinaryExpr(left Expr, op TokenType, right Expr) *BinaryExpr {
	return &BinaryExpr{left, op, right}
}

func NewUnaryExpr(op TokenType, e Expr) *UnaryExpr {
	return &UnaryExpr{op, e}
}

func NewBasicLit(t TokenType, v string) *BasicLit {
	return &BasicLit{t, v}
}

func NewParenExpr(e Expr) *ParenExpr {
	return &ParenExpr{e}
}

func NewAssignStmt(id string, e Expr) *AssignStmt {
	return &AssignStmt{id, e}
}

func NewReadStmt(v []string) *ReadStmt {
	return &ReadStmt{v}
}

func NewWriteStmt(v []string) *WriteStmt {
	return &WriteStmt{v}
}

func NewCondStmt(cond Expr, s ...[]Stmt) *CondStmt {
	if len(s) == 2 {
		return &CondStmt{cond, s[0], s[1]}
	}

	return &CondStmt{Cond: cond, If: s[0]}
}

func NewLoopStmt(cond Expr, s []Stmt) *LoopStmt {
	return &LoopStmt{cond, s}
}

func NewParser(s []Token) *Parser {
	return &Parser{source: s}
}

func (p *Parser) Parse() Node {
	return NewRootNode(p.series())
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
	return p.peek().Type == EOF
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
	p.consume(SEMICOLON)

	t := p.peek().Type
	for !p.isAtEnd() && !(t == END || t == ELSE || t == FI) {
		s = append(s, p.statement())
		p.consume(SEMICOLON)
		t = p.peek().Type
	}

	return s
}

func (p *Parser) statement() Stmt {
	switch p.peek().Type {
	case IDENTIFIER:
		return p.assignStmt()
	case READ:
		return p.readStmt()
	case WRITE:
		return p.writeStmt()
	case IF:
		return p.condStmt()
	default:
		return p.loopStmt()
	}
}

func (p *Parser) assignStmt() Stmt {
	id := p.nextToken()
	p.consume(COLON_EQUAL)

	return NewAssignStmt(id.Lexeme, p.expression())
}

func (p *Parser) readStmt() Stmt {
	p.consume(READ)

	return NewReadStmt(p.varList())
}

func (p *Parser) writeStmt() Stmt {
	p.consume(WRITE)

	return NewWriteStmt(p.varList())
}

func (p *Parser) condStmt() Stmt {
	p.consume(IF)
	cond := p.logical()
	p.consume(THEN)
	left := p.series()

	var right []Stmt
	if p.match(ELSE) {
		right = p.series()
	}
	p.consume(FI)

	if right != nil {
		return NewCondStmt(cond, left, right)
	}

	return NewCondStmt(cond, left)
}

func (p *Parser) loopStmt() Stmt {
	p.consume(WHILE)
	cond := p.logical()
	p.consume(DO)
	series := p.series()
	p.consume(END)

	return NewLoopStmt(cond, series)
}

func (p *Parser) varList() []string {
	vars := make([]string, 0)
	t := p.nextToken()
	vars = append(vars, t.Lexeme)

	for p.match(COMMMA) {
		vars = append(vars, p.nextToken().Lexeme)
	}

	return vars
}

func (p *Parser) logical() Expr {
	return p.logicalOr()
}

func (p *Parser) logicalOr() Expr {
	e := p.logicalAnd()

	for p.match(OR) {
		op := p.previous().Type
		right := p.logicalAnd()
		e = NewBinaryExpr(e, op, right)
	}

	return e
}

func (p *Parser) logicalAnd() Expr {
	e := p.logicalNot()

	for p.match(AND) {
		op := p.previous().Type
		right := p.logicalNot()
		e = NewBinaryExpr(e, op, right)
	}

	return e
}

func (p *Parser) logicalNot() Expr {
	if p.match(NOT) {
		return NewUnaryExpr(p.previous().Type, p.logical())
	}

	return p.logicalRelation()
}

func (p *Parser) logicalRelation() Expr {
	t := p.peek()
	if t.Type == TRUE || t.Type == FALSE || t.Type == LEFT_PAREN {
		return p.logicalElem()
	} else {
		return NewBinaryExpr(p.expression(), p.nextToken().Type, p.expression())
	}
}

func (p *Parser) logicalElem() Expr {
	t := p.nextToken()

	if t.Type == TRUE || t.Type == FALSE {
		return NewBasicLit(t.Type, t.Lexeme)
	} else {
		e := p.logical()
		p.consume(RIGHT_PAREN)
		return NewParenExpr(e)
	}
}

func (p *Parser) expression() Expr {
	exp := p.term()

	for p.match(PLUS) || p.match(MINUS) {
		op := p.previous().Type
		right := p.term()
		exp = NewBinaryExpr(exp, op, right)
	}

	return exp
}

func (p *Parser) term() Expr {
	elem := p.element()

	for p.match(STAR) || p.match(BACKSLASH) {
		op := p.previous().Type
		right := p.element()
		elem = NewBinaryExpr(elem, op, right)
	}

	return elem
}

func (p *Parser) element() Expr {
	t := p.nextToken()
	switch t.Type {
	case NUMBER, IDENTIFIER:
		return NewBasicLit(t.Type, t.Lexeme)
	default:
		elem := NewParenExpr(p.expression())
		p.consume(RIGHT_PAREN)
		return elem
	}
}
