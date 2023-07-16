package token

//go:generate stringer -type=TokenType

type TokenType int

const (
	NOT_FOUND TokenType = iota - 1
	ILLEGAL
	WHITESPACE
	COMMENT
	EOF

	literals_begin
	IDENT
	INT
	FLOAT
	BOOL_TRUE
	BOOL_FALSE
	STRING
	literals_end

	keywords_begin
	STRUCT
	TYPE
	LET
	SET
	TRUE
	FALSE
	NIL
	keywords_end

	symbols_begin
	LBRACKET
	RBRACKET
	LBRACE
	RBRACE
	LPAREN
	RPAREN
	PERIOD
	COMMA
	SEMICOLON
	COLON
	EQ
	ARROW
	symbols_end
)

var tokens = [...]string{
	ILLEGAL:    "ILLEGAL",
	WHITESPACE: "WHITESPACE",
	COMMENT:    "COMMENT",
	EOF:        "EOF",

	IDENT:      "IDENT",
	INT:        "INT",
	FLOAT:      "FLOAT",
	BOOL_TRUE:  "TRUE",
	BOOL_FALSE: "FALSE",
	STRING:     "STRING",

	STRUCT: "struct",
	TYPE:   "type",
	LET:    "let",
	SET:    "set",
	TRUE:   "true",
	FALSE:  "false",
	NIL:    "nil",

	LBRACKET:  "[",
	RBRACKET:  "]",
	LBRACE:    "{",
	RBRACE:    "}",
	LPAREN:    "(",
	RPAREN:    ")",
	PERIOD:    ".",
	COMMA:     ",",
	SEMICOLON: ";",
	COLON:     ":",
	EQ:        "=",
	ARROW:     "->",
}

var keywords map[string]TokenType

func init() {
	keywords = make(map[string]TokenType)
	for i := keywords_begin + 1; i < keywords_end; i++ {
		keywords[tokens[i]] = i
	}
}

func IsSymbol(r rune) bool {
	for i := symbols_begin + 1; i < symbols_end; i++ {
		if string(r) == tokens[i] {
			return true
		}
	}
	return false
}
func GetSymbol(r rune) TokenType {
	for i := symbols_begin + 1; i < symbols_end; i++ {
		if string(r) == tokens[i] {
			return i
		}
	}
	return NOT_FOUND
}

func IsKeyword(s string) bool {
	for i := keywords_begin + 1; i < keywords_end; i++ {
		if s == tokens[i] {
			return true
		}
	}
	return false
}
func GetKeyword(s string) TokenType {
	for i := keywords_begin + 1; i < keywords_end; i++ {
		if s == tokens[i] {
			return i
		}
	}
	return NOT_FOUND
}
