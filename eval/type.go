package eval

import (
	"fmt"
	"strings"
)

type TypeName struct {
	Name string
	Vars []TypeName
}

func (tn TypeName) String() string {
	name := tn.Name
	vars := []string{}
	for _, v := range tn.Vars {
		vars = append(vars, v.String())
	}
	if len(vars) == 0 {
		return name
	} else {
		return fmt.Sprintf("%s[%s]", name, strings.Join(vars, ","))
	}
}

type StructType struct {
	Fields map[string]TypeName
	Vars   []TypeName // positional typeargs
}

func (s StructType) String() string {
	fs := []string{}
	for id, tn := range s.Fields {
		fs = append(fs, fmt.Sprintf("%s %s", id, tn.String()))
	}
	vars := ""
	if len(s.Vars) > 0 {
		varNames := []string{}
		for _, v := range s.Vars {
			varNames = append(varNames, v.String())
		}
		vars = fmt.Sprintf("[%s]", strings.Join(varNames, ", "))
	}
	fields := strings.Join(fs, "; ")
	if len(s.Fields) == 0 {
		fields = ""
	}
	return fmt.Sprintf("struct%s{%s}", vars, fields)
}

func (s StructType) Copy() (o StructType) {
	o.Fields = make(map[string]TypeName)
	o.Vars = make([]TypeName, len(s.Vars))

	for f, tn := range s.Fields {
		o.Fields[f] = tn
	}
	copy(o.Vars, s.Vars)
	return o
}
