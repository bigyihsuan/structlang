# notes

- An esolang that is only structs
- Composition/embedding only
- Even ints are structs
- Struct defs are also structs... (?)
- <https://chat.stackexchange.com/transcript/message/63796738#63796738>

## types/structs

### builtin primitive types

everything is a struct. everything has fields for whatever reason.

- `int` (64 bit)
- `float` (64 bit)
- `bool`
- `string` (chars are 1-length strings)
- `nil` (only value of `nil` is `nil`)

included fields:

- `v`: itself
- `name`: the name of the type, as a string
- `len`: the "length" of the type. 0 for int, float, bool, nil.

### struct

- `struct`

### generics

- `struct[T]` (struct with type parameter)
- `either[T,U]` (builtin)

structs contain a list of fields

## making new types

```go
type number = int;
type point = struct { x,y float };
type list[T] = struct[T]{v T; next either[T,nil]}
```

## lexing info

- int: `[0-9]+`
- float: `[0-9]+\.[0-9]*`
- bool: `true|false`
- string: `".+"` except when escaped with `\`
- nil: `nil`
