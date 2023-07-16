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
}

func (s StructType) String() string {
	fs := []string{}
	for id, tn := range s.Fields {
		fs = append(fs, fmt.Sprintf("%s %s", id, tn.String()))
	}
	fields := strings.Join(fs, "; ")
	if len(s.Fields) == 0 {
		fields = ""
	}
	return fmt.Sprintf("struct{%s}", fields)
}
