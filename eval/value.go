package eval

type Value struct {
	V any
}

type StructValue struct {
	// TypeParams map[Identifier]StructType
	Fields map[Identifier]Value
}
