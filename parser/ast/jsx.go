package ast

import "yawp/parser/file"

type (
	JSXNode interface {
		_jsxNode()
	}

	JSXInterpolationNode struct {
		Start      file.Idx
		End        file.Idx
		Expression Expression
	}

	JSXElement struct {
		Start      file.Idx
		End        file.Idx
		Tag        Expression
		Attributes []JSXAttribute
		Children   []JSXNode
	}

	JSXFragment struct {
		Start    file.Idx
		End      file.Idx
		Children []JSXNode
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
		End        file.Idx
		Expression Expression
	}
)

func (*JSXInterpolationNode) _jsxNode() {}
func (*JSXElement) _jsxNode()           {}
func (*JSXFragment) _jsxNode()          {}

func (*JSXNamedAttribute) _jsxAttribute()  {}
func (*JSXSpreadAttribute) _jsxAttribute() {}

func (*JSXElement) _expressionNode()  {}
func (*JSXFragment) _expressionNode() {}

func (j *JSXElement) StartAt() file.Idx  { return j.Start }
func (j *JSXFragment) StartAt() file.Idx { return j.Start }

func (j *JSXElement) EndAt() file.Idx  { return j.End }
func (j *JSXFragment) EndAt() file.Idx { return j.End }
