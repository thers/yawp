package builtins

import (
	"yawp/parser/ast"
)

const emptyString = ""

var Arguments = &ast.Identifier{
	Ref:  &ast.Ref{
		Name: "arguments",
		Kind: ast.RBuiltin,
		Mangled: true,
	},
	Name: "arguments",
}

var SlicedArrayRest = &ast.Identifier{
	Ref: &ast.Ref{
		Name: "@slicedArrayRest",
		Kind: ast.RBuiltin,
	},
	Name: "@slicedArrayRest",
}

var ObjectRest = &ast.Identifier{
	Ref: &ast.Ref{
		Name: "@objectRest",
		Kind: ast.RBuiltin,
	},
	Name: "@objectRest",
}
