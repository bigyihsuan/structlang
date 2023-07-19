package lexer

import (
	"unicode"

	"github.com/bigyihsuan/structlang/token"
)

type Lexer struct {
	offset, line, column int
	src                  []rune
	lexeme               string
}

func NewLexer(src any) (*Lexer, error) {
	var s []rune

	switch src := src.(type) {
	case []byte:
		s = []rune(string(src))
	case string:
		s = []rune(src)
	case []rune:
		s = src
	}

	return &Lexer{
		offset: 0,
		line:   1,
		column: 1,
		src:    s,
		lexeme: "",
	}, nil
}

func (l *Lexer) Lex() token.Token {
	if l.offset >= len(l.src) {
		return token.NewToken(token.EOF, "EOF", l.offset, l.line, l.column)
	}
	offset, line, column := -1, -1, -1
	// skip whitespace
	for l.offset < len(l.src) {
		if isWhitespace(l.currentRune()) {
			if l.currentRune() == '\n' {
				l.nextLine()
			} else {
				l.nextCol()
			}
		} else {
			break
		}
	}
	if l.offset == len(l.src) {
		return token.NewToken(token.EOF, "", offset, line, column)
	}

	// set starting spans
	if offset == -1 {
		offset = l.offset
		line = l.line
		column = l.column
	}

	/* check for tokens */
	// comment
	if l.currentRune() == '/' {
		if next := l.nextRune(); next != -1 && next == '/' {
			l.comment()
			comment := l.resetLexeme()
			return token.NewToken(token.COMMENT, comment, offset, line, column)
		}
	}
	// arrow
	if l.currentRune() == '-' && l.nextRune() == '>' {
		l.addCurrent()
		l.addCurrent()
		lexeme := l.resetLexeme()
		return token.NewToken(token.ARROW, lexeme, offset, line, column)
	}
	// gteq, lteq,
	if l.currentRune() == '>' && l.nextRune() == '=' {
		l.addCurrent()
		l.addCurrent()
		lexeme := l.resetLexeme()
		return token.NewToken(token.GTEQ, lexeme, offset, line, column)
	}
	if l.currentRune() == '<' && l.nextRune() == '=' {
		l.addCurrent()
		l.addCurrent()
		lexeme := l.resetLexeme()
		return token.NewToken(token.LTEQ, lexeme, offset, line, column)
	}
	// int and float
	if isDigit(l.currentRune()) {
		numToken := l.intOrFloat()
		lexeme := l.resetLexeme()
		return token.NewToken(numToken, lexeme, offset, line, column)
	}
	// string
	if l.currentRune() == '"' {
		l.string()
		lexeme := l.resetLexeme()
		return token.NewToken(token.STRING, lexeme, offset, line, column)
	}
	// keywords and idents
	if isIdentOrKeywordChar(l.currentRune()) {
		l.addWhile(isIdentOrKeywordChar)
		lexeme := l.resetLexeme()
		if token.IsKeyword(lexeme) {
			kw := token.GetKeyword(lexeme)
			return token.NewToken(kw, lexeme, offset, line, column)
		} else {
			return token.NewToken(token.IDENT, lexeme, offset, line, column)
		}
	}
	// single characters
	symbolTokenType := token.GetSymbol(l.currentRune())
	if symbolTokenType != token.NOT_FOUND {
		l.addCurrent()
		lexeme := l.resetLexeme()
		return token.NewToken(symbolTokenType, lexeme, offset, line, column)
	}

	return token.NewToken(token.ILLEGAL, string(l.currentRune()), offset, line, column)
}

func (l *Lexer) comment() {
	l.addWhile(func(r rune) bool { return r != '\n' || l.offset >= len(l.src) })
	l.nextLine()
}
func (l *Lexer) intOrFloat() token.TokenType {
	l.addWhile(isDigit)
	if l.currentRune() != '.' {
		return token.INT
	}
	l.addCurrent()
	l.addWhile(isDigit)
	return token.FLOAT
}
func (l *Lexer) string() {
	l.nextCol() // ignore first quote
	for l.currentRune() != '"' {
		if l.currentRune() == '\\' {
			// escaped character
			l.addCurrent() // backslash
			l.addCurrent() // escaped char
		} else {
			l.addCurrent()
		}
	}
	l.nextCol() // ignore last quote
}

func (l Lexer) currentRune() rune { return l.src[l.offset] }
func (l Lexer) nextRune() rune {
	if l.offset >= len(l.src) {
		return -1
	}
	return l.src[l.offset+1]
}

func (l *Lexer) nextLine() {
	l.offset++
	l.line++
	l.column = 1
}
func (l *Lexer) nextCol() {
	l.offset++
	l.column++
}
func (l *Lexer) resetLexeme() string {
	lexeme := l.lexeme
	l.lexeme = ""
	return lexeme
}

// add the current character to the lexeme, and go to the next column.
func (l *Lexer) addCurrent() {
	l.lexeme += string(l.currentRune())
	l.nextCol()
}

// add the current character to the lexeme if some condition is met.
// returns the result of the predicate
func (l *Lexer) addWhile(condition func(r rune) bool) bool {
	ok := condition(l.currentRune())
	for ; ok; ok = condition(l.currentRune()) {
		l.addCurrent()
	}
	return ok
}

func isWhitespace(r rune) bool         { return unicode.IsSpace(r) }
func isDigit(r rune) bool              { return unicode.IsDigit(r) }
func isIdentOrKeywordChar(r rune) bool { return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' }

func ClearComments(tokens []token.Token) (out []token.Token) {
	for _, tok := range tokens {
		if tok.Type() != token.COMMENT {
			out = append(out, tok)
		}
	}
	return out
}
