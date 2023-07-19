package eval

import "fmt"

var builtinFuncs = map[string]func(vs ...Value){
	"print": func(vs ...Value) {
		if len(vs) == 0 {
			fmt.Print()
		}
		for _, v := range vs {
			fmt.Print(v.printString())
		}
	},
	"println": func(vs ...Value) {
		if len(vs) == 0 {
			fmt.Println()
		}
		for _, v := range vs {
			fmt.Println(v.printString())
		}
	},
}

func BuiltinFuncs() map[string]func(...Value) { return builtinFuncs }
