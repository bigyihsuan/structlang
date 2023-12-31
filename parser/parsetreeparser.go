package parser

import (
	"errors"
	"fmt"

	"github.com/bigyihsuan/structlang/token"
	"github.com/bigyihsuan/structlang/trees/parsetree"
	"github.com/bigyihsuan/structlang/trees/precedence"
	"github.com/bigyihsuan/structlang/util"
)

type ParseTreeParser struct {
	tokens    []token.Token
	idx       int
	prefixOps map[token.TokenType]PrefixParselet
	infixOps  map[token.TokenType]InfixParselet
}

func NewParser(tokens []token.Token) ParseTreeParser {
	prefixOps := make(map[token.TokenType]PrefixParselet)
	prefixOps = prefix(prefixOps, token.PLUS, precedence.PREFIX)
	prefixOps = prefix(prefixOps, token.MINUS, precedence.PREFIX)
	prefixOps = prefix(prefixOps, token.NOT, precedence.PREFIX)

	infixOps := make(map[token.TokenType]InfixParselet)
	infixOps = infixLeft(infixOps, token.PLUS, precedence.SUM)
	infixOps = infixLeft(infixOps, token.MINUS, precedence.SUM)
	infixOps = infixLeft(infixOps, token.STAR, precedence.PRODUCT)
	infixOps = infixLeft(infixOps, token.SLASH, precedence.PRODUCT)
	infixOps = infixLeft(infixOps, token.GT, precedence.COMPARISON)
	infixOps = infixLeft(infixOps, token.GTEQ, precedence.COMPARISON)
	infixOps = infixLeft(infixOps, token.LT, precedence.COMPARISON)
	infixOps = infixLeft(infixOps, token.LTEQ, precedence.COMPARISON)
	infixOps = infixLeft(infixOps, token.EQ, precedence.COMPARISON)
	infixOps = infixLeft(infixOps, token.AND, precedence.LOGICAL)
	infixOps = infixLeft(infixOps, token.OR, precedence.LOGICAL)

	// literals
	for _, primitive := range token.Primitives() {
		registerPrefix(prefixOps, primitive, LiteralParselet{})
	}
	registerPrefix(prefixOps, token.IDENT, IdentParselet{})
	registerPrefix(prefixOps, token.LPAREN, GroupingParselet{})
	registerInfix(infixOps, token.LPAREN, CallParselet{})
	registerPrefix(prefixOps, token.FUNC, FuncDefParselet{})

	return ParseTreeParser{
		tokens:    tokens,
		idx:       0,
		prefixOps: prefixOps,
		infixOps:  infixOps,
	}
}

func prefix(prefixOps map[token.TokenType]PrefixParselet, tt token.TokenType, prec precedence.Precedence) map[token.TokenType]PrefixParselet {
	registerPrefix(prefixOps, tt, PrefixOperator{prec})
	return prefixOps
}
func infixLeft(infixOps map[token.TokenType]InfixParselet, tt token.TokenType, prec precedence.Precedence) map[token.TokenType]InfixParselet {
	registerInfix(infixOps, tt, BinaryOperator{prec, false})
	return infixOps
}
func infixRight(infixOps map[token.TokenType]InfixParselet, tt token.TokenType, prec precedence.Precedence) map[token.TokenType]InfixParselet {
	registerInfix(infixOps, tt, BinaryOperator{prec, true})
	return infixOps
}
func registerPrefix(prefixOps map[token.TokenType]PrefixParselet, tt token.TokenType, pp PrefixParselet) {
	prefixOps[tt] = pp
}
func registerInfix(prefixOps map[token.TokenType]InfixParselet, tt token.TokenType, ip InfixParselet) {
	prefixOps[tt] = ip
}

func (p ParseTreeParser) Precedence() (precedence.Precedence, error) {
	tok, err := p.peekNextToken()
	if err != nil {
		return 0, err
	}
	infix, hasInfix := p.infixOps[tok.Type()]
	if hasInfix {
		return infix.Precedence(), nil
	} else {
		return precedence.BOTTOM, nil
	}
}

func (p ParseTreeParser) hasMoreTokens() bool {
	return p.idx < len(p.tokens)
}
func (p *ParseTreeParser) getNextToken() (tok *token.Token, err error) {
	if !p.hasMoreTokens() {
		return tok, errors.New("out of tokens")
	}
	tok = &p.tokens[p.idx]
	p.idx++
	return tok, nil
}
func (p ParseTreeParser) peekNextToken() (tok *token.Token, err error) {
	if !p.hasMoreTokens() {
		return tok, errors.New("out of tokens")
	}
	tok = &p.tokens[p.idx]
	return tok, nil
}
func (p *ParseTreeParser) putBackToken() {
	p.idx--
}

func (p *ParseTreeParser) expectGet(tt token.TokenType) (*token.Token, error) {
	tok, err := p.getNextToken()
	if err != nil {
		return tok, err
	}
	if tok.Type() != tt {
		return tok, fmt.Errorf("expected token `%s`, got `%s` at `%v`", tt, tok.Type(), tok.Position())
	}
	return tok, nil
}
func (p *ParseTreeParser) expectGetAny(tts ...token.TokenType) (*token.Token, error) {
	tok, err := p.getNextToken()
	if err != nil {
		return tok, err
	}
	for _, tt := range tts {
		if tok.Type() == tt {
			return tok, nil
		}
	}
	return tok, fmt.Errorf("expected any token %s, got `%s` at `%v`", tts, tok.Type(), tok.Position())
}

func (p ParseTreeParser) nextTokenIs(tt token.TokenType) (bool, error) {
	next, err := p.peekNextToken()
	if err != nil {
		return false, err
	}
	return next.Type() == tt, nil
}
func (p ParseTreeParser) nextTokenIsAny(tts ...token.TokenType) (bool, error) {
	next, err := p.peekNextToken()
	if err != nil {
		return false, err
	}
	for _, tt := range tts {
		if next.Type() == tt {
			return true, nil
		}
	}
	return false, nil
}

func (p *ParseTreeParser) Parse() (stmts []parsetree.Stmt, errs error) {
	for p.hasMoreTokens() {
		if hasRbrace, err := p.nextTokenIs(token.RBRACE); err != nil {
			return stmts, err
		} else if hasRbrace {
			return stmts, nil
		}
		s, e := p.Stmt()
		if e != nil {
			return stmts, errors.Join(errors.New("in parse tree parser"), e)
		}
		stmts = append(stmts, s)
	}
	return stmts, errs
}
func (p *ParseTreeParser) Stmt() (stmt parsetree.Stmt, errs error) {
	stmterr := errors.New("in stmt")
	kw, err := p.peekNextToken()
	if err != nil {
		return stmt, errors.Join(stmterr, errors.New("missing keyword token"), err)
	}
	switch kw.Type() {
	case token.TYPE:
		td, err := p.TypeDef()
		if err != nil {
			return td, errors.Join(stmterr, errors.New("expected typedef with kw `type`"), err)
		}
		return td, nil
	case token.LET:
		vd, err := p.VarDef()
		if err != nil {
			return vd, errors.Join(stmterr, errors.New("expected vardef with kw `let`"), err)
		}
		return vd, nil
	case token.SET:
		vs, err := p.VarSet()
		if err != nil {
			return vs, errors.Join(stmterr, errors.New("expected vardef with kw `let`"), err)
		}
		return vs, nil
	case token.RETURN:
		rs, err := p.ReturnStmt()
		if err != nil {
			return rs, errors.Join(stmterr, errors.New("expected return with kw `return`"), err)
		}
		return rs, nil
	default:
		expr, err := p.ExprStmt()
		if err != nil {
			return expr, errors.Join(stmterr, err)
		}
		return expr, nil
		// return nil, errors.Join(stmterr, fmt.Errorf("unknown for Stmt: type=`%s` lexeme=`%s`", kw.Type(), kw.Lexeme()))
	}
}

func (p *ParseTreeParser) VarDef() (vd parsetree.VarDef, errs error) {
	vderr := errors.New("in vardef")
	letkw, err := p.expectGet(token.LET)
	if err != nil {
		return vd, errors.Join(vderr, err)
	}
	lvalue, err := p.Lvalue()
	if err != nil {
		return vd, errors.Join(vderr, errors.New("expected lvalue"), err)
	}
	eq, err := p.expectGet(token.EQ)
	if err != nil {
		return vd, errors.Join(vderr, err)
	}
	rvalue, err := p.Expr(precedence.BOTTOM)
	if err != nil {
		return vd, errors.Join(vderr, errors.New("expected rvalue"), err)
	}
	sc, err := p.expectGet(token.SEMICOLON)
	if err != nil {
		return vd, errors.Join(vderr, err)
	}
	return parsetree.VarDef{LetKw: *letkw, Lvalue: lvalue, Eq: *eq, Rvalue: rvalue, Sc: *sc}, nil
}

func (p *ParseTreeParser) VarSet() (vs parsetree.VarSet, errs error) {
	vserr := errors.New("in varset")
	setkw, err := p.expectGet(token.SET)
	if err != nil {
		return vs, errors.Join(vserr, err)
	}
	lvalue, err := p.Lvalue()
	if err != nil {
		return vs, errors.Join(vserr, errors.New("expected lvalue"), err)
	}
	eq, err := p.expectGet(token.EQ)
	if err != nil {
		return vs, errors.Join(vserr, err)
	}
	rvalue, err := p.Expr(precedence.BOTTOM)
	if err != nil {
		return vs, errors.Join(vserr, errors.New("expected rvalue"), err)
	}
	sc, err := p.expectGet(token.SEMICOLON)
	if err != nil {
		return vs, errors.Join(vserr, err)
	}
	return parsetree.VarSet{SetKw: *setkw, Lvalue: lvalue, Eq: *eq, Rvalue: rvalue, Sc: *sc}, nil
}

func (p *ParseTreeParser) Lvalue() (lv parsetree.Lvalue, err error) {
	lverr := errors.New("in lvalue")
	// check for ident first
	ident, err := p.expectGet(token.IDENT)
	if err != nil {
		return lv, errors.Join(lverr, err)
	}
	if isArrow, err := p.nextTokenIs(token.ARROW); err != nil {
		return lv, errors.Join(lverr, err)
	} else if isArrow {
		p.putBackToken() // put back ident
		return p.FieldAccess()
	} else {
		// single ident
		return parsetree.Ident{Name: *ident}, nil
	}
}

func (p *ParseTreeParser) TypeDef() (td parsetree.TypeDef, errs error) {
	tderr := errors.New("in typedef")
	type_, err := p.expectGet(token.TYPE)
	if err != nil {
		return td, errors.Join(tderr, err)
	}
	typename, err := p.Type()
	if err != nil {
		return td, errors.Join(tderr, errors.New("expected typename"), err)
	}
	eq, err := p.expectGet(token.EQ)
	if err != nil {
		return td, errors.Join(tderr, err)
	}
	structDef, err := p.StructDef()
	if err != nil {
		return td, errors.Join(tderr, errors.New("expected structdef"), err)
	}
	sc, err := p.expectGet(token.SEMICOLON)
	if err != nil {
		return td, errors.Join(tderr, err)
	}

	return parsetree.TypeDef{TypeKw: *type_, TypeName: typename, Eq: *eq, StructDef: structDef, Sc: *sc}, nil
}

func (p *ParseTreeParser) Type() (ty parsetree.Type, errs error) {
	tyerr := errors.New("in type")
	typename, err := p.Ident()
	if err != nil {
		return ty, errors.Join(tyerr, err)
	}
	typevars, err := p.TypeVars()
	if err != nil {
		return ty, errors.Join(tyerr, errors.New("expected typevars"), err)
	}
	return parsetree.Type{TypeName: typename, TypeVars: typevars}, nil
}

func (p *ParseTreeParser) TypeVars() (tvs *parsetree.TypeVars, errs error) {
	tvserr := errors.New("in typevars")
	if peeked, err := p.peekNextToken(); err != nil {
		return tvs, errors.Join(tvserr, err)
	} else if peeked.Type() != token.LBRACKET {
		return nil, nil
	}
	lbracket, err := p.expectGet(token.LBRACKET)
	if err != nil {
		return tvs, errors.Join(tvserr, err)
	}
	typevars, err := p.TypeVarParams()
	if err != nil {
		return tvs, errors.Join(tvserr, errors.New("expected typevar params"), err)
	}
	rbracket, err := p.expectGet(token.RBRACKET)
	if err != nil {
		return tvs, errors.Join(tvserr, err)
	}
	return &parsetree.TypeVars{Lbracket: *lbracket, TypeVars: typevars, Rbracket: *rbracket}, nil
}

func (p *ParseTreeParser) TypeVarParams() (tv parsetree.SeparatedList[parsetree.Type, token.Token], errs error) {
	tvperr := errors.New("in typevar params")
	for {
		if peeked, err := p.peekNextToken(); err != nil {
			return tv, errors.Join(tvperr, err)
		} else if peeked.Type() != token.IDENT && peeked.Type() != token.NIL {
			// 0 eles, or with trailing sep
			return tv, nil
		}
		typename, err := p.Type()
		if err != nil {
			return tv, errors.Join(tvperr, errors.New("expected typename"), err)
		}
		if peeked, err := p.peekNextToken(); err != nil {
			return tv, errors.Join(tvperr, err)
		} else if tt := peeked.Type(); tt == token.RBRACE || tt != token.COMMA {
			// with no trailing sep
			tv = append(tv, util.Pair[parsetree.Type, *token.Token]{First: typename, Last: nil})
			return
		}
		comma, err := p.expectGet(token.COMMA)
		if err != nil {
			return tv, errors.Join(tvperr, err)
		}
		tv = append(tv, util.Pair[parsetree.Type, *token.Token]{First: typename, Last: comma})
	}
}

func (p *ParseTreeParser) StructDef() (st parsetree.StructDef, errs error) {
	sderr := errors.New("in structdef")
	structKw, err := p.expectGet(token.STRUCT)
	if err != nil {
		return st, errors.Join(sderr, err)
	}
	typeVars, err := p.TypeVars()
	if err != nil {
		return st, errors.Join(sderr, errors.New("expected typevars"), err)
	}
	lbrace, err := p.expectGet(token.LBRACE)
	if err != nil {
		return st, errors.Join(sderr, err)
	}
	fields, err := p.StructFields()
	if err != nil {
		return st, errors.Join(sderr, errors.New("expected struct fields"), err)
	}
	rbrace, err := p.expectGet(token.RBRACE)
	if err != nil {
		return st, errors.Join(sderr, err)
	}
	return parsetree.StructDef{StructKw: *structKw, TypeVars: typeVars, Lbrace: *lbrace, Fields: fields, Rbrace: *rbrace}, nil
}

func (p *ParseTreeParser) StructFields() (f []parsetree.StructField, errs error) {
	sferr := errors.New("in struct fields")
	for {
		if peeked, err := p.peekNextToken(); err != nil {
			return f, errors.Join(sferr, err)
		} else if tt := peeked.Type(); tt == token.RBRACE {
			// exit when rbrace
			return f, nil
		}
		names, err := p.NameList()
		if err != nil {
			return f, errors.Join(sferr, errors.New("expected name list"), err)
		} else if len(names) < 1 {
			return f, errors.Join(sferr, errors.New("name list must be len > 0"))
		}
		typename, err := p.Type()
		if err != nil {
			return f, errors.Join(sferr, errors.New("expected typename"), err)
		}
		if peeked, err := p.peekNextToken(); err != nil {
			return f, errors.Join(sferr, err)
		} else if tt := peeked.Type(); tt == token.RBRACE {
			// exit when names-typename pair, but no sc
			f = append(f, parsetree.StructField{Names: names, Type: typename, Sc: nil})
			return f, nil
		}
		sc, err := p.expectGet(token.SEMICOLON)
		if err != nil {
			return f, errors.Join(sferr, err)
		}
		f = append(f, parsetree.StructField{Names: names, Type: typename, Sc: sc})
		// no trailing scs allowed
	}
}

func (p *ParseTreeParser) NameList() (names parsetree.SeparatedList[parsetree.Ident, token.Token], errs error) {
	nlerr := errors.New("in namelist")
	for {
		if peeked, err := p.peekNextToken(); err != nil {
			return names, errors.Join(nlerr, err)
		} else if tt := peeked.Type(); tt != token.IDENT {
			// exit when not ident
			return names, nil
		}
		name, err := p.Ident()
		if err != nil {
			return names, errors.Join(nlerr, err)
		}
		if peeked, err := p.peekNextToken(); err != nil {
			return names, errors.Join(nlerr, err)
		} else if tt := peeked.Type(); tt != token.COMMA {
			// exit when ident but no comma
			names = append(names, util.Pair[parsetree.Ident, *token.Token]{First: name, Last: nil})
			return names, nil
		}
		comma, err := p.expectGet(token.COMMA)
		if err != nil {
			return names, errors.Join(nlerr, err)
		}
		names = append(names, util.Pair[parsetree.Ident, *token.Token]{First: name, Last: comma})
		// no trailing comma allowed
	}
}

func (p *ParseTreeParser) ReturnStmt() (stmt parsetree.ReturnStmt, err error) {
	rserr := errors.New("in returnstmt")
	returnKw, err := p.expectGet(token.RETURN)
	if err != nil {
		return stmt, errors.Join(rserr, err)
	}
	var expr parsetree.Expr
	if hasSc, err := p.nextTokenIs(token.SEMICOLON); err != nil {
		return stmt, errors.Join(rserr, err)
	} else if !hasSc {
		expr, err = p.Expr(precedence.BOTTOM)
		if err != nil {
			return stmt, err
		}
	}
	sc, err := p.expectGet(token.SEMICOLON)
	if err != nil {
		return stmt, errors.Join(rserr, err)
	}
	return parsetree.ReturnStmt{
		ReturnKw: *returnKw,
		Expr:     expr,
		Sc:       *sc,
	}, nil
}

func (p *ParseTreeParser) ExprStmt() (stmt parsetree.ExprStmt, err error) {
	eserr := errors.New("in exprstmt")
	e, err := p.Expr(precedence.BOTTOM)
	if err != nil {
		return stmt, errors.Join(eserr, err)
	}
	sc, err := p.expectGet(token.SEMICOLON)
	if err != nil {
		return stmt, errors.Join(eserr, err)
	}
	return parsetree.ExprStmt{Expr: e, Sc: *sc}, nil
}

func (p *ParseTreeParser) Expr(precedence precedence.Precedence) (expr parsetree.Expr, err error) {
	exprerr := errors.New("in expr")
	tok, err := p.getNextToken()
	if err != nil {
		return expr, errors.Join(exprerr, err)
	}
	pp, isPrefix := p.prefixOps[tok.Type()]
	if !isPrefix {
		return expr, errors.Join(exprerr, fmt.Errorf("could not get prefix parslet for token `%s`", tok))
	}
	expr, err = pp.Parse(p, *tok)
	if err != nil {
		return expr, errors.Join(exprerr, err)
	}

	nextPrecedence, err := p.Precedence()
	if err != nil {
		return expr, errors.Join(exprerr, err)
	}
	for precedence < nextPrecedence {
		op, err := p.getNextToken()
		if err != nil {
			return expr, errors.Join(exprerr, err)
		}
		infix, hasInfix := p.infixOps[op.Type()]
		if !hasInfix {
			return expr, errors.Join(exprerr, fmt.Errorf("could not get infix parslet for token `%s`", tok))
		}
		expr, err = infix.Parse(p, expr, *op)
		if err != nil {
			return expr, errors.Join(exprerr, err)
		}
		nextPrecedence, err = p.Precedence()
		if err != nil {
			return expr, errors.Join(exprerr, err)
		}
	}
	return expr, nil
}

func (p *ParseTreeParser) Literal() (expr parsetree.Literal, err error) {
	literr := errors.New("in literal")
	tok, err := p.getNextToken()
	if err != nil {
		return expr, errors.Join(literr, err)
	}
	return parsetree.Literal{Token: *tok}, nil
}

func (p *ParseTreeParser) IdentOrStructLiteralOrFieldAccess() (expr parsetree.Expr, err error) {
	islerr := errors.New("in ident/struct literal/field access")
	ident, err := p.expectGet(token.IDENT)
	if err != nil {
		return expr, errors.Join(islerr, err)
	}
	if hasStructLiteral, err := p.nextTokenIsAny(token.LBRACE, token.LBRACKET); err != nil {
		return expr, errors.Join(islerr, err)
	} else if hasStructLiteral {
		p.putBackToken()
		sl, err := p.StructLiteral()
		if err != nil {
			return expr, errors.Join(islerr, errors.New("expected struct literal with `{`"), err)
		}
		return sl, nil
	}
	if hasFieldAccess, err := p.nextTokenIs(token.ARROW); err != nil {
		return expr, errors.Join(islerr, err)
	} else if hasFieldAccess {
		p.putBackToken()
		fa, err := p.FieldAccess()
		if err != nil {
			return expr, errors.Join(islerr, errors.New("expected field access with `->`"), err)
		}
		return fa, nil
	}
	return parsetree.Ident{Name: *ident}, nil
}

func (p *ParseTreeParser) Ident() (i parsetree.Ident, err error) {
	ierr := errors.New("in ident")
	name, err := p.expectGetAny(token.IDENT, token.NIL)
	if err != nil {
		return i, errors.Join(ierr, err)
	}
	return parsetree.Ident{Name: *name}, nil
}

func (p *ParseTreeParser) StructLiteral() (sl parsetree.StructLiteral, err error) {
	slerr := errors.New("in struct literal")
	typename, err := p.Type()
	if err != nil {
		return sl, errors.Join(slerr, errors.New("expected type"), err)
	}
	lbrace, err := p.expectGet(token.LBRACE)
	if err != nil {
		return sl, errors.Join(slerr, err)
	}
	fields, err := p.StructLiteralFields()
	if err != nil {
		return sl, errors.Join(slerr, errors.New("expected struct literal fields"), err)
	}
	rbrace, err := p.expectGet(token.RBRACE)
	if err != nil {
		return sl, errors.Join(slerr, err)
	}
	return parsetree.StructLiteral{TypeName: typename, Lbrace: *lbrace, Fields: fields, Rbrace: *rbrace}, nil
}

func (p *ParseTreeParser) StructLiteralFields() (slfs parsetree.SeparatedList[parsetree.StructLiteralField, token.Token], err error) {
	slfserr := errors.New("in struct literal fields")
	for {
		if tok, err := p.peekNextToken(); err != nil {
			return slfs, errors.Join(slfserr, err)
		} else if tok.Type() == token.RBRACE {
			// 0 or many fields
			return slfs, nil
		}
		fieldName, err := p.Ident()
		if err != nil {
			return slfs, errors.Join(slfserr, err)
		}
		colon, err := p.expectGet(token.COLON)
		if err != nil {
			return slfs, errors.Join(slfserr, err)
		}
		value, err := p.Expr(precedence.BOTTOM)
		if err != nil {
			return slfs, errors.Join(slfserr, errors.New("expected expr"), err)
		}
		field := parsetree.StructLiteralField{FieldName: fieldName, Colon: *colon, Value: value}
		if tok, err := p.peekNextToken(); err != nil {
			return slfs, errors.Join(slfserr, err)
		} else if tok.Type() == token.RBRACE {
			// 1 field
			pair := util.Pair[parsetree.StructLiteralField, *token.Token]{First: field, Last: nil}
			slfs = append(slfs, pair)
			return slfs, nil
		}
		comma, err := p.expectGet(token.COMMA)
		if err != nil {
			return slfs, errors.Join(slfserr, err)
		}
		pair := util.Pair[parsetree.StructLiteralField, *token.Token]{First: field, Last: comma}
		slfs = append(slfs, pair)
	}
}

func (p *ParseTreeParser) FieldAccess() (fa parsetree.Lvalue, err error) {
	faerr := errors.New("in field access")
	fa, err = p.Ident()
	if err != nil {
		return fa, errors.Join(faerr, errors.New("expected lvalue with `->`"), err)
	}
	for {
		if hasArrow, err := p.nextTokenIs(token.ARROW); err != nil {
			return fa, errors.Join(faerr, err)
		} else if !hasArrow {
			break
		}
		arrow, err := p.expectGet(token.ARROW)
		if err != nil {
			return fa, errors.Join(faerr, err)
		}
		field, err := p.expectGet(token.IDENT)
		if err != nil {
			return fa, errors.Join(faerr, errors.New("expected identifier with `->`"), err)
		}
		fa = parsetree.FieldAccess{
			Lvalue: fa,
			Arrow:  *arrow,
			Field:  parsetree.Ident{Name: *field},
		}
	}
	return fa, err
}
