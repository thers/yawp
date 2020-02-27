package parser

import (
	"yawp/parser/ast"
	"yawp/parser/file"
	"yawp/parser/token"
)

func (p *Parser) parseObjectPropertyFromShorthand(idx file.Idx, literal string) *ast.Property {
	if p.is(token.COMMA) || matchIdentifier.MatchString(p.literal) {
		p.consumePossible(token.COMMA)

		return &ast.Property{
			Key:  literal,
			Kind: "value",
			Value: &ast.Identifier{
				Name: literal,
				Idx:  idx,
			},
		}
	}

	p.unexpectedToken()
	return nil
}

func (p *Parser) parseObjectProperty() *ast.Property {
	if p.is(token.IDENTIFIER) {
		shouldConsumeNext := true

		// Maybe setter or getter
		if p.literal == "get" || p.literal == "set" {
			idx, literal := p.idx, p.literal
			p.next()

			// Setter or getter
			if p.is(token.IDENTIFIER) {
				propertyKey := p.literal
				p.next()
				parameterList := p.parseFunctionParameterList()

				node := &ast.FunctionLiteral{
					Function:   idx,
					Parameters: parameterList,
				}
				p.parseFunctionBlock(node)
				p.consumePossible(token.COMMA)

				return &ast.Property{
					Key:   propertyKey,
					Kind:  literal,
					Value: node,
				}
			}

			shouldConsumeNext = false
		}

		idx, propertyKey := p.idx, p.literal
		if shouldConsumeNext {
			p.next()
		}

		// Object function shorthand
		if p.is(token.LEFT_PARENTHESIS) {
			parameterList := p.parseFunctionParameterList()
			node := &ast.FunctionLiteral{
				Function:   idx,
				Parameters: parameterList,
			}

			p.parseFunctionBlock(node)
			node.Source = p.slice(idx, node.Body.Idx1())

			p.consumePossible(token.COMMA)

			return &ast.Property{
				Key:   propertyKey,
				Kind:  "function",
				Value: node,
			}
		}

		if p.is(token.COLON) {
			p.consumeExpected(token.COLON)

			return &ast.Property{
				Key:   propertyKey,
				Kind:  "value",
				Value: p.parseAssignmentExpression(),
			}
		}

		return p.parseObjectPropertyFromShorthand(idx, propertyKey)
	} else if p.is(token.NUMBER) {
		key := p.literal
		_, err := parseNumberLiteral(p.literal)

		if err != nil {
			key = ""
			p.error(p.idx, err.Error())
		}

		p.consumeExpected(token.COLON)

		return &ast.Property{
			Key:   key,
			Kind:  "value",
			Value: p.parseAssignmentExpression(),
		}
	} else if p.is(token.STRING) {
		key, err := parseStringLiteral(p.literal[1 : len(p.literal)-1])
		if err != nil {
			p.error(p.idx, err.Error())
		}

		p.consumeExpected(token.COLON)

		return &ast.Property{
			Key:   key,
			Kind:  "value",
			Value: p.parseAssignmentExpression(),
		}
	}

	p.next()
	return p.parseObjectPropertyFromShorthand(p.idx, p.literal)
}

func (p *Parser) parseObjectLiteral() ast.Expression {
	var value []*ast.Property

	idx0 := p.consumeExpected(token.LEFT_BRACE)
	for !p.is(token.RIGHT_BRACE) && !p.is(token.EOF) {
		property := p.parseObjectProperty()
		value = append(value, property)

		if p.is(token.RIGHT_BRACE) {
			break
		}
	}
	idx1 := p.consumeExpected(token.RIGHT_BRACE)

	return &ast.ObjectLiteral{
		LeftBrace:  idx0,
		RightBrace: idx1,
		Value:      value,
	}
}
