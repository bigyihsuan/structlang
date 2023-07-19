package eval

type Value interface {
	Get(field string) Value
	TypeName() TypeName
	Unwrap() any
}

type StructValue struct {
	Name       string
	TypeParams map[string]TypeName
	Fields     map[string]Value
}

func NewStructValue(fields map[string]Value) (sv StructValue) {
	sv.Fields = make(map[string]Value)
	for name, value := range fields {
		sv.Fields[name] = value
	}
	return
}
func NewStructValueFromType(st StructType, typeParams map[string]TypeName, fields map[string]Value, typeName string) (sv StructValue) {
	sv.TypeParams = typeParams
	sv.Name = typeName
	sv.Fields = make(map[string]Value)
	for name := range st.Fields {
		var v Value
		sv.Fields[name] = v
	}
	for name := range sv.Fields {
		sv.Fields[name] = fields[name]
	}
	return
}

func (sv StructValue) Get(field string) Value {
	return sv.Fields[field]
}

func (sv StructValue) TypeName() TypeName {
	return TypeName{Name: sv.Name}
}
func (sv StructValue) Unwrap() any {
	// TODO: what is this unwrapped?
	return sv.Fields
}
