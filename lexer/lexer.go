package lexer

import (
	"strconv"
	"unicode"

	"structlang/token"
)

type Lexer struct {
	offset    int
	lookAhead int
	line      int
	column    int
	filename  string
	src       []rune
	ch        rune
}

func New(src any, filename string) *Lexer {
	var s []rune

	switch src := src.(type) {
	case []byte:
		s = []rune(string(src))
	case string:
		s = []rune(src)
	case []rune:
		s = src
	}

	l := &Lexer{
		offset:    0,
		lookAhead: 0,
		line:      1,
		column:    1,
		src:       s,
		filename:  filename,
	}
	l.readChar()

	return l
}

func (l *Lexer) GetNextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '[':
		tok = token.New(token.LBRACKET, string(l.ch), l.filename, l.offset, l.line, l.column)
	case ']':
		tok = token.New(token.RBRACKET, string(l.ch), l.filename, l.offset, l.line, l.column)
	case '{':
		tok = token.New(token.LBRACE, string(l.ch), l.filename, l.offset, l.line, l.column)
	case '}':
		tok = token.New(token.RBRACE, string(l.ch), l.filename, l.offset, l.line, l.column)
	case '(':
		tok = token.New(token.LPAREN, string(l.ch), l.filename, l.offset, l.line, l.column)
	case ')':
		tok = token.New(token.RPAREN, string(l.ch), l.filename, l.offset, l.line, l.column)
	case '.':
		tok = token.New(token.PERIOD, string(l.ch), l.filename, l.offset, l.line, l.column)
	case ',':
		tok = token.New(token.COMMA, string(l.ch), l.filename, l.offset, l.line, l.column)
	case ';':
		tok = token.New(token.SEMICOLON, string(l.ch), l.filename, l.offset, l.line, l.column)
	case ':':
		tok = token.New(token.COLON, string(l.ch), l.filename, l.offset, l.line, l.column)
	case '=':
		lexeme := string(l.ch)
		if l.peekChar() == '=' {
			l.readChar()
			lexeme += string(l.ch)
			tok = token.New(token.DEQ, lexeme, l.filename, l.offset, l.line, l.column)
		} else {
			tok = token.New(token.EQ, string(l.ch), l.filename, l.offset, l.line, l.column)
		}
	case '+':
		tok = token.New(token.PLUS, string(l.ch), l.filename, l.offset, l.line, l.column)
	case '-':
		lexeme := string(l.ch)
		if l.peekChar() == '>' {
			l.readChar()
			lexeme += string(l.ch)
			tok = token.New(token.ARROW, lexeme, l.filename, l.offset, l.line, l.column)
		} else {
			tok = token.New(token.MINUS, string(l.ch), l.filename, l.offset, l.line, l.column)
		}
	case '*':
		tok = token.New(token.STAR, string(l.ch), l.filename, l.offset, l.line, l.column)
	case '/':
		tok = token.New(token.SLASH, string(l.ch), l.filename, l.offset, l.line, l.column)
	case '>':
		lexeme := string(l.ch)
		if l.peekChar() == '=' {
			l.readChar()
			lexeme += string(l.ch)
			tok = token.New(token.GTEQ, lexeme, l.filename, l.offset, l.line, l.column)
		} else {
			tok = token.New(token.GT, string(l.ch), l.filename, l.offset, l.line, l.column)
		}
	case '<':
		lexeme := string(l.ch)
		if l.peekChar() == '=' {
			l.readChar()
			lexeme += string(l.ch)
			tok = token.New(token.LTEQ, lexeme, l.filename, l.offset, l.line, l.column)
		} else {
			tok = token.New(token.LT, string(l.ch), l.filename, l.offset, l.line, l.column)
		}
	case '"':
		tok.Type = token.STRING
		tok.Lexeme = l.string()
	case 0:
		tok.Type = token.EOF
		tok.Lexeme = ""
	default:
		if isLetter(l.ch) {
			tok.Lexeme = l.ident()
			tok.Type = token.GetKeywordOrIdent(tok.Lexeme)
			return tok
		} else if isDigit(l.ch) {
			tok.Lexeme, tok.Type = l.intOrFloat()
			return tok
		} else {
			tok = token.New(token.ILLEGAL, string(l.ch), l.filename, l.offset, l.line, l.column)
		}
	}
	l.readChar()
	return tok
}

func (l *Lexer) readChar() {
	if l.lookAhead >= len(l.src) {
		l.ch = 0
	} else {
		l.ch = l.src[l.lookAhead]
	}
	l.offset = l.lookAhead
	l.column++
	l.lookAhead++
	if l.ch == '\n' {
		l.line++
		l.column = 1
	}
}
func (l *Lexer) peekChar() rune {
	if l.lookAhead >= len(l.src) {
		return 0
	} else {
		return l.src[l.lookAhead]
	}
}

func (l *Lexer) skipWhitespace() {
	for isWhitespace(l.ch) {
		l.readChar()
		if l.ch == '\n' {
			l.line++
			l.column = 1
		}
	}
}

func (l *Lexer) ident() string {
	startPosition := l.offset
	for isLetter(l.ch) {
		l.readChar()
	}
	return string(l.src[startPosition:l.offset])
}
func (l *Lexer) intOrFloat() (string, token.TokenType) {
	tt := token.INT
	startPosition := l.offset
	for isDigit(l.ch) {
		l.readChar()
	}
	if l.ch == '.' {
		tt = token.FLOAT
		l.readChar()
		for isDigit(l.ch) {
			l.readChar()
		}
	}
	return string(l.src[startPosition:l.offset]), tt
}
func (l *Lexer) string() string {
	startPosition := l.offset + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	s, _ := strconv.Unquote("\"" + string(l.src[startPosition:l.offset]) + "\"")
	return s
}

func isLetter(ch rune) bool     { return unicode.IsLetter(ch) }
func isWhitespace(ch rune) bool { return unicode.IsSpace(ch) }
func isDigit(ch rune) bool      { return unicode.IsDigit(ch) }
