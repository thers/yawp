package transpiler

import (
	"strconv"
	"yawp/builtins"
	"yawp/parser/ast"
)

func createArraySlice(array ast.IExpr, index int) *ast.CallExpression {
	return &ast.CallExpression{
		Callee: builtins.SlicedArrayRest,
		ArgumentList: []ast.IExpr{
			array,
			&ast.NumberLiteral{
				Literal: strconv.Itoa(index),
			},
		},
	}
}
