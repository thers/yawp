package builtins

import (
	"yawp/parser/ast"
)

const emptyString = ""

var Arguments = &ast.Identifier{
	LegacyRef:  &ast.SymbolRef{
		Name:    "arguments",
		Type:    ast.SRBuiltin,
		Mangled: true,
	},
	Name: "arguments",
}

var SlicedArrayRest = &ast.Identifier{
	LegacyRef: &ast.SymbolRef{
		Name: "@slicedArrayRest",
		Type: ast.SRBuiltin,
	},
	Name: "@slicedArrayRest",
}

var ObjectRest = &ast.Identifier{
	LegacyRef: &ast.SymbolRef{
		Name: "@objectRest",
		Type: ast.SRBuiltin,
	},
	Name: "@objectRest",
}
