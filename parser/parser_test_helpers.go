package parser

import (
	"structlang/ast"
	"testing"
)

func testProgram(t *testing.T, program *ast.Program, stmtCount int) bool {
	if len(program.Stmts) != stmtCount {
		t.Errorf("program stmt count incorrect: want=%d got=%d", stmtCount, len(program.Stmts))
		return false
	}
	return true
}

func wrongTypeFatal(t *testing.T, kind string, value any) {
	wrongType(t, kind, value)
	t.FailNow()
}

func wrongType(t *testing.T, kind string, value any) {
	t.Errorf("not a %s: got=%T", kind, value)
}
