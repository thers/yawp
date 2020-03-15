package ast

import "yawp/parser/file"

type (
	TemplateExpression struct {
		Start         file.Idx
		End           file.Idx
		Strings       []string
		Substitutions []Expression
	}

	TaggedTemplateExpression struct {
		Tag      Expression
		Template *TemplateExpression
	}
)

func (*TemplateExpression) _expressionNode()       {}
func (*TaggedTemplateExpression) _expressionNode() {}

func (t *TemplateExpression) StartAt() file.Idx       { return t.Start }
func (t *TaggedTemplateExpression) StartAt() file.Idx { return t.Tag.StartAt() }

func (t *TemplateExpression) EndAt() file.Idx       { return t.End }
func (t *TaggedTemplateExpression) EndAt() file.Idx { return t.Template.EndAt() }
