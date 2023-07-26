// Code generated by "stringer -type=TokenType"; DO NOT EDIT.

package token

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[NOT_FOUND - -1]
	_ = x[ILLEGAL-0]
	_ = x[WHITESPACE-1]
	_ = x[COMMENT-2]
	_ = x[EOF-3]
	_ = x[IDENT-4]
	_ = x[literals_begin-5]
	_ = x[INT-6]
	_ = x[FLOAT-7]
	_ = x[STRING-8]
	_ = x[literals_end-9]
	_ = x[keywords_begin-10]
	_ = x[STRUCT-11]
	_ = x[TYPE-12]
	_ = x[LET-13]
	_ = x[SET-14]
	_ = x[TRUE-15]
	_ = x[FALSE-16]
	_ = x[NIL-17]
	_ = x[AND-18]
	_ = x[OR-19]
	_ = x[NOT-20]
	_ = x[FUNC-21]
	_ = x[RETURN-22]
	_ = x[keywords_end-23]
	_ = x[symbols_begin-24]
	_ = x[LBRACKET-25]
	_ = x[RBRACKET-26]
	_ = x[LBRACE-27]
	_ = x[RBRACE-28]
	_ = x[LPAREN-29]
	_ = x[RPAREN-30]
	_ = x[PERIOD-31]
	_ = x[COMMA-32]
	_ = x[SEMICOLON-33]
	_ = x[COLON-34]
	_ = x[EQ-35]
	_ = x[ARROW-36]
	_ = x[PLUS-37]
	_ = x[MINUS-38]
	_ = x[STAR-39]
	_ = x[SLASH-40]
	_ = x[GT-41]
	_ = x[LT-42]
	_ = x[GTEQ-43]
	_ = x[LTEQ-44]
	_ = x[symbols_end-45]
}

const _TokenType_name = "NOT_FOUNDILLEGALWHITESPACECOMMENTEOFIDENTliterals_beginINTFLOATSTRINGliterals_endkeywords_beginSTRUCTTYPELETSETTRUEFALSENILANDORNOTFUNCRETURNkeywords_endsymbols_beginLBRACKETRBRACKETLBRACERBRACELPARENRPARENPERIODCOMMASEMICOLONCOLONEQARROWPLUSMINUSSTARSLASHGTLTGTEQLTEQsymbols_end"

var _TokenType_index = [...]uint16{0, 9, 16, 26, 33, 36, 41, 55, 58, 63, 69, 81, 95, 101, 105, 108, 111, 115, 120, 123, 126, 128, 131, 135, 141, 153, 166, 174, 182, 188, 194, 200, 206, 212, 217, 226, 231, 233, 238, 242, 247, 251, 256, 258, 260, 264, 268, 279}

func (i TokenType) String() string {
	i -= -1
	if i < 0 || i >= TokenType(len(_TokenType_index)-1) {
		return "TokenType(" + strconv.FormatInt(int64(i+-1), 10) + ")"
	}
	return _TokenType_name[_TokenType_index[i]:_TokenType_index[i+1]]
}
