package parser

import (
	"errors"
	"fmt"

	"github.com/bigyihsuan/structlang/token"
	"github.com/bigyihsuan/structlang/trees/parsetree"
	"github.com/bigyihsuan/structlang/util"
)

type ParseTreeParser struct {
	tokens []token.Token
	idx    int
}

func NewParser(tokens []token.Token) ParseTreeParser {
	return ParseTreeParser{
		tokens: tokens,
		idx:    0,
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

func (p *ParseTreeParser) Parse() (stmts []parsetree.Stmt, errs error) {
	for p.hasMoreTokens() {
		s, e := p.Stmt()
		if e != nil {
			return stmts, e
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
	default:
		return nil, errors.Join(stmterr, fmt.Errorf("unknown for Stmt: type=`%s` lexeme=`%s`", kw.Type(), kw.Lexeme()))
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
	typename, err := p.expectGetAny(token.IDENT, token.NIL)
	if err != nil {
		return ty, errors.Join(tyerr, err)
	}
	typevars, err := p.TypeVars()
	if err != nil {
		return ty, errors.Join(tyerr, errors.New("expected typevars"), err)
	}
	return parsetree.Type{TypeName: *typename, TypeVars: typevars}, nil
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
func (p *ParseTreeParser) NameList() (names parsetree.SeparatedList[token.Token, token.Token], errs error) {
	nlerr := errors.New("in namelist")
	for {
		if peeked, err := p.peekNextToken(); err != nil {
			return names, errors.Join(nlerr, err)
		} else if tt := peeked.Type(); tt != token.IDENT {
			// exit when not ident
			return names, nil
		}
		name, err := p.expectGet(token.IDENT)
		if err != nil {
			return names, errors.Join(nlerr, err)
		}
		if peeked, err := p.peekNextToken(); err != nil {
			return names, errors.Join(nlerr, err)
		} else if tt := peeked.Type(); tt != token.COMMA {
			// exit when ident but no comma
			names = append(names, util.Pair[token.Token, *token.Token]{First: *name, Last: nil})
			return names, nil
		}
		comma, err := p.expectGet(token.COMMA)
		if err != nil {
			return names, errors.Join(nlerr, err)
		}
		names = append(names, util.Pair[token.Token, *token.Token]{First: *name, Last: comma})
		// no trailing comma allowed
	}
}
