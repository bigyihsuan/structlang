// Code generated by "stringer -type=TokenType"; DO NOT EDIT.

package token

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[NOT_FOUND-1]
	_ = x[ILLEGAL-2]
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
	_ = x[DEQ-43]
	_ = x[GTEQ-44]
	_ = x[LTEQ-45]
	_ = x[symbols_end-46]
}

const _TokenType_name = "NOT_FOUNDILLEGALEOFIDENTliterals_beginINTFLOATSTRINGliterals_endkeywords_beginSTRUCTTYPELETSETTRUEFALSENILANDORNOTFUNCRETURNkeywords_endsymbols_beginLBRACKETRBRACKETLBRACERBRACELPARENRPARENPERIODCOMMASEMICOLONCOLONEQARROWPLUSMINUSSTARSLASHGTLTDEQGTEQLTEQsymbols_end"

var _TokenType_index = [...]uint16{0, 9, 16, 19, 24, 38, 41, 46, 52, 64, 78, 84, 88, 91, 94, 98, 103, 106, 109, 111, 114, 118, 124, 136, 149, 157, 165, 171, 177, 183, 189, 195, 200, 209, 214, 216, 221, 225, 230, 234, 239, 241, 243, 246, 250, 254, 265}

func (i TokenType) String() string {
	i -= 1
	if i < 0 || i >= TokenType(len(_TokenType_index)-1) {
		return "TokenType(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _TokenType_name[_TokenType_index[i]:_TokenType_index[i+1]]
}
