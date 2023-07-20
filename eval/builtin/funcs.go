package builtin

import (
	"fmt"

	. "github.com/bigyihsuan/structlang/eval"
)

var builtinFuncs = map[string]func(vs ...Value) Value{
	"print":   print_,
	"println": println_,
}

func BuiltinFuncs() map[string]func(...Value) Value { return builtinFuncs }

func print_(vs ...Value) Value {
	if len(vs) == 0 {
		fmt.Print()
	}
	for _, v := range vs {
		fmt.Print(v.PrintString())
	}
	return NewNil()
}

func println_(vs ...Value) Value {
	if len(vs) == 0 {
		fmt.Println()
	}
	for _, v := range vs {
		fmt.Println(v.PrintString())
	}
	return NewNil()
}
