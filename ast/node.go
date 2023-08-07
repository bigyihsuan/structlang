package ast

import (
	"fmt"
	"structlang/token"
)

type Node interface {
	BoundingTokens() (first token.Token, last token.Token)
	fmt.Stringer
}
