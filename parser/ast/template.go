package ast

import "yawp/parser/file"

type (
	TemplateExpression struct {
		Start         file.Loc
		End           file.Loc
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

func (t *TemplateExpression) StartAt() file.Loc       { return t.Start }
func (t *TaggedTemplateExpression) StartAt() file.Loc { return t.Tag.StartAt() }

func (t *TemplateExpression) EndAt() file.Loc       { return t.End }
func (t *TaggedTemplateExpression) EndAt() file.Loc { return t.Template.EndAt() }
