package eval

import (
	"fmt"
	"strings"
)

type TypeName struct {
	Name Identifier
	Vars []TypeName
}

func (tn TypeName) String() string {
	name := tn.Name.String()
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
	TypeParams map[Identifier]TypeName
	Fields     map[Identifier]TypeName
}

func (s StructType) String() string {
	tps := []string{}
	for _, tn := range s.TypeParams {
		tps = append(tps, tn.String())
	}
	typeParams := fmt.Sprintf("[%s]", strings.Join(tps, ", "))
	if len(s.TypeParams) == 0 {
		typeParams = ""
	}

	fs := []string{}
	for id, tn := range s.Fields {
		fs = append(fs, fmt.Sprintf("%s %s", id.String(), tn.String()))
	}
	fields := strings.Join(fs, "; ")
	if len(s.Fields) == 0 {
		fields = ""
	}
	return fmt.Sprintf("struct%s{%s}", typeParams, fields)
}
