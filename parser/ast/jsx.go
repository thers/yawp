package ast

import "yawp/parser/file"

type (
	JSXElementName struct {
		Expression IExpr
		StringName string
	}

	JSXAttribute interface {
		_jsxAttribute()
	}

	JSXNamedAttribute struct {
		Name  *Identifier
		Value IExpr
	}

	JSXSpreadAttribute struct {
		Start      file.Idx
		Expression IExpr
	}

	JSXChild interface {
		_jsxChild()
	}

	JSXText struct {
		Node
		Text string
	}

	JSXChildExpression struct {
		IExpr
	}
)

func (*JSXChildExpression) _jsxChild() {}
func (*JSXText) _jsxChild()            {}
func (*JSXElement) _jsxChild()         {}
func (*JSXFragment) _jsxChild()        {}

func (*JSXNamedAttribute) _jsxAttribute()  {}
func (*JSXSpreadAttribute) _jsxAttribute() {}
