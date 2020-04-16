package ast

import "yawp/parser/file"

type (
	JSXElement struct {
		Loc        *file.Loc
		Name       *JSXElementName
		Attributes []JSXAttribute
		Children   []JSXChild
	}

	JSXElementName struct {
		Expression Expression
		StringName string
	}

	JSXNamespacedName struct {
		Loc       *file.Loc
		Namespace string
		Name      string
	}

	JSXFragment struct {
		Loc      *file.Loc
		Children []JSXChild
	}

	JSXAttribute interface {
		_jsxAttribute()
	}

	JSXNamedAttribute struct {
		Name  *Identifier
		Value Expression
	}

	JSXSpreadAttribute struct {
		Start      file.Idx
		Expression Expression
	}

	JSXChild interface {
		_jsxChild()
	}

	JSXText struct {
		Loc  *file.Loc
		Text string
	}

	JSXChildExpression struct {
		Expression
	}
)

func (*JSXChildExpression) _jsxChild() {}
func (*JSXText) _jsxChild()            {}
func (*JSXElement) _jsxChild()         {}
func (*JSXFragment) _jsxChild()        {}

func (*JSXNamedAttribute) _jsxAttribute()  {}
func (*JSXSpreadAttribute) _jsxAttribute() {}

func (*JSXElement) _expressionNode()        {}
func (*JSXFragment) _expressionNode()       {}
func (*JSXNamespacedName) _expressionNode() {}

func (j *JSXElement) GetLoc() *file.Loc        { return j.Loc }
func (j *JSXFragment) GetLoc() *file.Loc       { return j.Loc }
func (j *JSXNamespacedName) GetLoc() *file.Loc { return j.Loc }
