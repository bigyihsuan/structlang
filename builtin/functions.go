package builtin

import (
	"fmt"
	"strings"

	"github.com/bigyihsuan/structlang/env"
	"github.com/bigyihsuan/structlang/trees/ast"
	"github.com/bigyihsuan/structlang/util"
	. "github.com/bigyihsuan/structlang/value"
)

type Func struct { // function name
	Args     []util.Pair[string, TypeName] // function arguments; oredered map of variable names to their types
	Body     []ast.Stmt                    // the actual code of the function
	Env      *env.Env                      // the env the parameter variables live in
	IsReturn bool
}

func (f Func) Get(field string) Value {
	return NewNil()
}
func (f Func) TypeName() TypeName {
	args := make([]string, len(f.Args))
	for _, a := range f.Args {
		args = append(args, a.Last.String())
	}
	return TypeName{Name: fmt.Sprintf("func(%s)", strings.Join(args, ", "))}
}
func (f Func) Unwrap() any {
	return NewNil()
}
func (f Func) PrintString() string {
	args := []string{}
	for _, p := range f.Args {
		name := p.First
		ty := p.Last
		args = append(args, fmt.Sprintf("%s %s", name, ty))
	}
	return fmt.Sprintf("func(%s)", strings.Join(args, ", "))
}
func (f Func) Return(isReturn bool) Value {
	f.IsReturn = isReturn
	return f
}

func (f Func) Call(evaluator Eval, args ...Value) (Value, error) {
	if len(args) != len(f.Args) {
		// TODO: better runtime error handling
		panic(fmt.Sprintf("incorrect numbers of arguments for func: got %d, want %d", len(args), len(f.Args)))
	}
	for i, argValue := range args {
		argName := f.Args[i].First
		argType := f.Args[i].Last
		if argValue.TypeName().Name != argType.Name {
			panic(fmt.Sprintf("incorrect argument types for func: got %s, want %s", argValue.TypeName(), argType))
		}
		f.Env.DefineVariable(argName, argValue)
	}

	// TODO: return values
	retVal, err := evaluator.Evaluate(f.Env, f.Body)
	return retVal, err
}
