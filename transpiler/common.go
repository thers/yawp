package transpiler

import (
	"strconv"
	"yawp/builtins"
	"yawp/parser/ast"
)

func createArraySlice(array ast.Expression, index int) *ast.CallExpression {
	return &ast.CallExpression{
		Callee: builtins.SlicedArrayRest,
		ArgumentList: []ast.Expression{
			array,
			&ast.NumberLiteral{
				Literal: strconv.Itoa(index),
			},
		},
	}
}
