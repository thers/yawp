package ast

import "yawp/parser/file"

type (
	TemplateExpression struct {
		Loc           *file.Loc
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

func (t *TemplateExpression) GetLoc() *file.Loc { return t.Loc }
func (t *TaggedTemplateExpression) GetLoc() *file.Loc {
	return t.Tag.GetLoc().Add(t.Template.GetLoc())
}
