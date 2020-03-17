package ast

import "yawp/parser/file"

type (
	JSXElement struct {
		Start      file.Idx
		End        file.Idx
		Name       *JSXElementName
		Attributes []JSXAttribute
		Children   []JSXChild
	}

	JSXElementName struct {
		Expression Expression
		StringName string
	}

	JSXNamespacedName struct {
		Start     file.Idx
		Namespace string
		Name      string
	}

	JSXFragment struct {
		Start    file.Idx
		End      file.Idx
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
		Start file.Idx
		End   file.Idx
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

func (j *JSXElement) StartAt() file.Idx        { return j.Start }
func (j *JSXFragment) StartAt() file.Idx       { return j.Start }
func (j *JSXNamespacedName) StartAt() file.Idx { return j.Start }

func (j *JSXElement) EndAt() file.Idx  { return j.End }
func (j *JSXFragment) EndAt() file.Idx { return j.End }
func (j *JSXNamespacedName) EndAt() file.Idx {
	return j.Start + file.Idx(len(j.Name)) + 1 + file.Idx(len(j.Namespace))
}
