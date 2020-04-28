package builtins

import (
	"yawp/parser/ast"
)

const emptyString = ""

var SlicedArrayRest = &ast.Identifier{
	Ref: &ast.Ref{
		Name: emptyString,
		Kind: ast.RBuiltin,
	},
	Name: emptyString,
}

var ObjectRest = &ast.Identifier{
	Ref: &ast.Ref{
		Name: emptyString,
		Kind: ast.RBuiltin,
	},
	Name: emptyString,
}
