package ast

import (
	"fmt"
	"structlang/token"
)

type Ident struct {
	Token token.Token
	Value string
}

func (expr *Ident) exprTag()   {}
func (expr *Ident) lvalueTag() {}
func (expr *Ident) BoundingTokens() (token.Token, token.Token) {
	return expr.Token, expr.Token
}
func (expr *Ident) String() string { return fmt.Sprintf("(%s)", expr.Value) }

type IntLit struct {
	Token token.Token
	Value int64
}

func (expr *IntLit) exprTag() {}
func (expr *IntLit) BoundingTokens() (token.Token, token.Token) {
	return expr.Token, expr.Token
}
func (expr *IntLit) String() string { return fmt.Sprintf(`("%s")`, expr.Token.Lexeme) }

type FloatLit struct {
	Token token.Token
	Value float64
}

func (expr *FloatLit) exprTag() {}
func (expr *FloatLit) BoundingTokens() (token.Token, token.Token) {
	return expr.Token, expr.Token
}
func (expr *FloatLit) String() string { return fmt.Sprintf("(%s)", expr.Token.Lexeme) }

type StrLit struct {
	Token token.Token
	Value string
}

func (expr *StrLit) exprTag() {}
func (expr *StrLit) BoundingTokens() (token.Token, token.Token) {
	return expr.Token, expr.Token
}
func (expr *StrLit) String() string { return fmt.Sprintf("(%s)", expr.Token.Lexeme) }

type BoolLit struct {
	Token token.Token
	Value bool
}

func (expr *BoolLit) exprTag() {}
func (expr *BoolLit) BoundingTokens() (token.Token, token.Token) {
	return expr.Token, expr.Token
}
func (expr *BoolLit) String() string { return fmt.Sprintf("(%s)", expr.Token.Lexeme) }
