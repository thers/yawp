package parser

import (
	"yawp/parser/ast"
	"yawp/parser/token"
)

func (p *Parser) parseTemplateExpression() *ast.TemplateExpression {
	exp := &ast.TemplateExpression{
		Strings:       make([]string, 0),
		Substitutions: make([]ast.Expression, 0),
	}

	exp.Start = p.idx

	currentString := ""

	// we're parsing template literal in chars mode mostly
	// parsing until we meet another `
	for p.chr != '`' {
		// escape, so the next chr will be just add to the current string part
		if p.chr == '\\' {
			// advance to the next chr after \
			p.read()

			// add current chr to string no matter what it is
			currentString += string(p.chr)

			// and advance to the next one
			p.read()
			continue
		} else if p.chr == '$' {
			// advance to the next chr
			p.read()

			// start of substitution
			if p.chr == '{' {
				// adding current string to strings
				exp.Strings = append(exp.Strings, currentString)
				currentString = ""

				// advance to the next chr as we're at { now
				p.read()
				// and scan next token as substitution will be parsed as expression
				// thus switching to "tokens" mode back from chars
				p.next()
				exp.Substitutions = append(exp.Substitutions, p.parseAssignmentExpression())

				// we're still in "tokens" mode, so check for next token instead of chr
				if !p.is(token.RIGHT_BRACE) {
					p.unexpectedToken()
					p.next()
					return nil
				}

				continue
			} else {
				currentString += "$"
				continue
			}
		}

		currentString += string(p.chr)
		p.read()
	}

	// reading past last ` so we can normally back to "tokens" mode
	p.read()
	exp.Strings = append(exp.Strings, currentString)
	exp.End = p.idx

	// advance to the next token
	p.next()

	return exp
}

func (p *Parser) parseTaggedTemplateExpression(tag ast.Expression) *ast.TaggedTemplateExpression {
	return &ast.TaggedTemplateExpression{
		Tag:      tag,
		Template: p.parseTemplateExpression(),
	}
}
