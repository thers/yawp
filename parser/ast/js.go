package ast

import "yawp/parser/file"

// Like ASM embeds in C this is to embed raw js
// for obvious reasons this code is only ES5 and won't be transpiled
type Js struct {
	Code string
}

func (*Js) GetLoc() *file.Loc { return nil }
func (*Js) GetNode() *Node    { return nil }
func (*Js) _expressionNode()  {}
func (*Js) _statementNode()   {}
