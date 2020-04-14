package ast

import "yawp/parser/file"

type (
	JSXElement struct {
		Start      file.Loc
		End        file.Loc
		Name       *JSXElementName
		Attributes []JSXAttribute
		Children   []JSXChild
	}

	JSXElementName struct {
		Expression Expression
		StringName string
	}

	JSXNamespacedName struct {
		Start     file.Loc
		Namespace string
		Name      string
	}

	JSXFragment struct {
		Start    file.Loc
		End      file.Loc
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
		Start      file.Loc
		Expression Expression
	}

	JSXChild interface {
		_jsxChild()
	}

	JSXText struct {
		Start file.Loc
		End   file.Loc
		Text  string
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

func (j *JSXElement) StartAt() file.Loc        { return j.Start }
func (j *JSXFragment) StartAt() file.Loc       { return j.Start }
func (j *JSXNamespacedName) StartAt() file.Loc { return j.Start }

func (j *JSXElement) EndAt() file.Loc  { return j.End }
func (j *JSXFragment) EndAt() file.Loc { return j.End }
func (j *JSXNamespacedName) EndAt() file.Loc {
	return j.Start + file.Loc(len(j.Name)) + 1 + file.Loc(len(j.Namespace))
}
