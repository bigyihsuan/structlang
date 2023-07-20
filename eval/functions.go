package eval

import (
	"fmt"
	"strings"

	"github.com/bigyihsuan/structlang/trees/ast"
	"github.com/bigyihsuan/structlang/util"
)

type Func struct {
	Name   string                        // function name
	Args   []util.Pair[string, TypeName] // function arguments; oredered map of variable names to their types
	Return TypeName                      // the return type of this function
	Body   []ast.Stmt                    // the actual code of the function
	Env    *Env                          // the env the parameter variables live in
}

func (f Func) Get(field string) Value {
	return NewNil()
}
func (f Func) TypeName() TypeName {
	args := make([]string, len(f.Args))
	for _, a := range f.Args {
		args = append(args, a.Last.String())
	}
	return TypeName{Name: fmt.Sprintf("func(%s)%s", strings.Join(args, ", "), f.Return.Name)}
}
func (f Func) Unwrap() any {
	return NewNil()
}

func (f Func) Call(evaluator Eval, args ...Value) Value {
	if len(args) != len(f.Args) {
		// TODO: better runtime error handling
		panic(fmt.Sprintf("incorrect numbers of arguments for func %s: got %d, want %d", f.Name, len(args), len(f.Args)))
	}
	for i, argValue := range args {
		argName := f.Args[i].First
		argType := f.Args[i].Last
		if argValue.TypeName().Name != argType.Name {
			panic(fmt.Sprintf("incorrect argument types for func %s: got %s, want %s", f.Name, argValue.TypeName(), argType))
		}
		f.Env.DefineVariable(argName, argValue)
	}

	// TODO: return values
	evaluator.Evaluate(f.Env, f.Body)

	return NewNil()
}
