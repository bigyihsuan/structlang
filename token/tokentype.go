package token

//go:generate stringer -type=TokenType

type TokenType int

const (
	NOT_FOUND TokenType = iota - 1
	ILLEGAL
	WHITESPACE
	COMMENT
	EOF

	IDENT
	literals_begin
	INT
	FLOAT
	// BOOL_TRUE
	// BOOL_FALSE
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
	// AND
	// OR
	// NOT
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
	PLUS
	MINUS
	STAR
	SLASH
	symbols_end
)

var tokens = [...]string{
	ILLEGAL:    "ILLEGAL",
	WHITESPACE: "WHITESPACE",
	COMMENT:    "COMMENT",
	EOF:        "EOF",

	IDENT: "IDENT",
	INT:   "INT",
	FLOAT: "FLOAT",
	// BOOL_TRUE:  "TRUE",
	// BOOL_FALSE: "FALSE",
	STRING: "STRING",

	STRUCT: "struct",
	TYPE:   "type",
	LET:    "let",
	SET:    "set",
	TRUE:   "true",
	FALSE:  "false",
	NIL:    "nil",
	// AND:    "and",
	// OR:     "or",
	// NOT:    "not",

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
	PLUS:      "+",
	MINUS:     "-",
	STAR:      "*",
	SLASH:     "/",
}

var primitives []TokenType

var keywords map[string]TokenType

func init() {
	keywords = make(map[string]TokenType)
	for i := keywords_begin + 1; i < keywords_end; i++ {
		keywords[tokens[i]] = i
	}

	for i := literals_begin + 1; i < literals_end; i++ {
		primitives = append(primitives, i)
	}
	primitives = append(primitives, TRUE)
	primitives = append(primitives, FALSE)
	primitives = append(primitives, NIL)
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

func IsLiteral(tt TokenType) bool {
	return (literals_begin < tt || tt < literals_end) || tt == NIL || tt == FALSE || tt == TRUE
}

func Primitives() []TokenType {
	return primitives
}
