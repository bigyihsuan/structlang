package ast

import (
	"fmt"
	"strings"
	"structlang/token"
)

type Program struct {
	Stmts []Stmt
}

func (p Program) BoundingTokens() (token.Token, token.Token) {
	s := token.New(token.ILLEGAL, "", "", -1, -1, -1)
	e := token.New(token.ILLEGAL, "", "", -1, -1, -1)
	if len(p.Stmts) > 0 {
		s, _ = p.Stmts[0].BoundingTokens()
		_, e = p.Stmts[len(p.Stmts)-1].BoundingTokens()
	}
	return s, e
}
func (p Program) String() string {
	var builder strings.Builder
	for _, s := range p.Stmts {
		builder.WriteString(s.String())
	}
	return builder.String()
}

type Stmt interface {
	stmtTag()
	Node
}

type ExprStmt struct {
	FirstToken, LastToken token.Token
	Expr                  Expr
}

func (s ExprStmt) stmtTag()                                   {}
func (s ExprStmt) BoundingTokens() (token.Token, token.Token) { return s.FirstToken, s.LastToken }
func (s ExprStmt) String() string {
	if s.Expr != nil {
		return s.Expr.String()
	}
	return ""
}

type LetStmt struct {
	FirstToken, LastToken token.Token
	Ident                 Lvalue
	Value                 Expr
}

func (stmt *LetStmt) stmtTag() {}
func (stmt *LetStmt) BoundingTokens() (token.Token, token.Token) {
	return stmt.FirstToken, stmt.LastToken
}
func (stmt *LetStmt) String() string {
	return fmt.Sprintf("(let %s = %s)", stmt.Ident.String(), stmt.Value.String())
}
