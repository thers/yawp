package parser

import (
	"yawp/parser/ast"
	"yawp/parser/file"
	"yawp/parser/token"
)

func (p *Parser) parseFlowTypeIdentifierRemainder(identifier *ast.Identifier) *ast.FlowIdentifier {
	qualifier := &ast.FlowIdentifier{
		Start: identifier.Start,
		Name:  identifier.Name,
	}

	for {
		if p.is(token.PERIOD) {
			p.next()

			identifier = p.parseIdentifier()
			qualifier = &ast.FlowIdentifier{
				Start:         identifier.Start,
				Name:          identifier.Name,
				Qualification: qualifier,
			}
		} else {
			return qualifier
		}
	}
}

func (p *Parser) parseFlowTypeIdentifier() *ast.FlowIdentifier {
	identifier := p.parseIdentifier()

	return p.parseFlowTypeIdentifierRemainder(identifier)
}

func (p *Parser) parseFlowTypeIdentifierIncludingKeywords() *ast.FlowIdentifier {
	identifier := p.parseIdentifierIncludingKeywords()

	if identifier == nil {
		return nil
	}

	return &ast.FlowIdentifier{
		Start: identifier.Start,
		Name:  identifier.Name,
	}
}

func (p *Parser) parseSimpleFlowType() ast.FlowType {
	start := p.loc

	switch p.token {
	case token.BOOLEAN:
		kind := p.literal
		p.next()

		if kind == "true" {
			return &ast.FlowTrueType{
				Start: start,
			}
		} else {
			return &ast.FlowFalseType{
				Start: start,
			}
		}
	case token.TYPE_BOOLEAN, token.TYPE_ANY, token.TYPE_STRING, token.TYPE_NUMBER, token.VOID, token.NULL, token.TYPE_MIXED:
		kind := p.token
		end := start + file.Loc(len(p.literal))
		p.next()

		return &ast.FlowPrimitiveType{
			Start: start,
			End:   end,
			Kind:  kind,
		}
	case token.STRING:
		str := p.literal[1 : len(p.literal)-1]
		p.next()

		return &ast.FlowStringLiteralType{
			Start:  start,
			String: str,
		}
	case token.NUMBER:
		number, err := parseNumberLiteral(p.literal)
		p.next()

		if err != nil {
			p.error(start, err.Error())
		} else {
			return &ast.FlowNumberLiteralType{
				Start:  start,
				Number: number,
			}
		}
	case token.IDENTIFIER:
		id := p.parseFlowTypeIdentifier()

		if p.isFlowTypeArgumentsStart() {
			return &ast.FlowGenericType{
				Name:          id,
				TypeArguments: p.parseFlowTypeArguments(),
			}
		} else {
			return id
		}
	case token.TYPEOF:
		p.next()

		return &ast.FlowTypeOfType{
			Start:      start,
			Identifier: p.parseFlowTypeIdentifier(),
		}
	case token.QUESTION_MARK:
		p.next()

		return &ast.FlowOptionalType{
			FlowType: p.parseFlowType(),
		}
	case token.MULTIPLY:
		p.next()

		return &ast.FlowExistentialType{
			Start: start,
		}
	case token.LEFT_BRACE:
		return p.parseFlowInexactObjectType()
	case token.TYPE_EXACT_OBJECT_START:
		return p.parseFlowExactObjectType()
	case token.LEFT_BRACKET:
		return p.parseFlowTupleType()
	case token.LESS:
		typeParameters := p.parseFlowTypeParameters()
		start := p.consumeExpected(token.LEFT_PARENTHESIS)
		parameters := p.parseFlowFunctionParameters()
		p.consumeExpected(token.RIGHT_PARENTHESIS)

		functionType := p.parseFlowFunctionRemainder(start, parameters)
		functionType.TypeParameters = typeParameters

		return functionType
	}

	return nil
}

func (p *Parser) parseFlowFunctionRemainder(start file.Loc, params []*ast.FlowFunctionParameter) *ast.FlowFunctionType {
	p.consumeExpected(token.ARROW)

	return &ast.FlowFunctionType{
		Start:      start,
		Parameters: params,
		ReturnType: p.parseFlowType(),
	}
}

func (p *Parser) parseFlowFunctionParameter() *ast.FlowFunctionParameter {
	if p.is(token.IDENTIFIER) {
		// possible var identifier
		identifier := p.parseIdentifier()

		if p.is(token.COLON) {
			p.next()

			return &ast.FlowFunctionParameter{
				Identifier: identifier,
				Type:       p.parseFlowType(),
			}
		}
	}

	return &ast.FlowFunctionParameter{
		Type: p.parseFlowType(),
	}
}

func (p *Parser) parseFlowFunctionParameters() []*ast.FlowFunctionParameter {
	parameters := make([]*ast.FlowFunctionParameter, 0)

	for p.until(token.RIGHT_PARENTHESIS) {
		parameters = append(parameters, p.parseFlowFunctionParameter())

		p.consumePossible(token.COMMA)
	}

	return parameters
}

func (p *Parser) parseFlowExpressionOrFunction() ast.FlowType {
	start := p.consumeExpected(token.LEFT_PARENTHESIS)

	var flowType ast.FlowType

	if !p.is(token.RIGHT_PARENTHESIS) {
		closeScope := p.openTypeScope()
		flowType = p.parseFlowType()
		closeScope()
	}

	if p.isAny(token.COMMA, token.COLON) {
		parameter := &ast.FlowFunctionParameter{
			Type: flowType,
		}

		if id, ok := flowType.(*ast.FlowIdentifier); p.is(token.COLON) && ok {
			p.next()
			parameter.Identifier = id.Identifier()
			parameter.Type = p.parseFlowType()
		} else {
			p.next()
		}

		// it's functions params

		parameters := []*ast.FlowFunctionParameter{
			parameter,
		}

		parameters = append(parameters, p.parseFlowFunctionParameters()...)

		p.consumeExpected(token.RIGHT_PARENTHESIS)

		return p.parseFlowFunctionRemainder(start, parameters)
	}

	if p.is(token.RIGHT_PARENTHESIS) {
		// exp or 0-1 params fn
		p.next()

		if p.is(token.ARROW) {
			parameters := make([]*ast.FlowFunctionParameter, 0)

			if flowType != nil {
				parameters = append(parameters, &ast.FlowFunctionParameter{
					Type: flowType,
				})
			}

			return p.parseFlowFunctionRemainder(start, parameters)
		}

		return flowType
	}

	return nil
}

func (p *Parser) parseFlowUnionType(beginning ast.FlowType) *ast.FlowUnionType {
	p.scope.allowUnionType = false

	defer func() {
		p.scope.allowUnionType = true
	}()

	// already parsed `type`
	p.consumeExpected(token.OR)

	union := &ast.FlowUnionType{
		Types: make([]ast.FlowType, 0),
	}

	if beginning == nil {
		union.Start = p.loc
	} else {
		// FIXME: add StartAt() to ast.FlowType
		union.Start = -1
		union.Types = append(union.Types, beginning)
	}

	// second element
	union.Types = append(union.Types, p.parseFlowType())

	for !p.is(token.EOF) && p.is(token.OR) {
		p.next()

		union.Types = append(union.Types, p.parseFlowType())
	}

	return union
}

func (p *Parser) parseFlowIntersectionType(beginning ast.FlowType) *ast.FlowIntersectionType {
	p.scope.allowIntersectionType = false

	defer func() {
		p.scope.allowIntersectionType = true
	}()

	// already parsed `type`
	p.consumeExpected(token.AND)

	intersection := &ast.FlowIntersectionType{
		Types: make([]ast.FlowType, 0),
	}

	if beginning == nil {
		intersection.Start = p.loc
	} else {
		// FIXME: add StartAt() to ast.FlowType
		intersection.Start = -1
		intersection.Types = append(intersection.Types, beginning)
	}

	// second element
	intersection.Types = append(intersection.Types, p.parseFlowType())

	for !p.is(token.EOF) && p.is(token.AND) {
		p.next()

		intersection.Types = append(intersection.Types, p.parseFlowType())
	}

	return intersection
}

func (p *Parser) parseFlowType() ast.FlowType {
	start := p.loc
	var flowType ast.FlowType

	// could be type expression enclosure or function args
	if p.is(token.LEFT_PARENTHESIS) {
		flowType = p.parseFlowExpressionOrFunction()
	} else {
		// if not, parse simple type
		flowType = p.parseSimpleFlowType()
	}

	// it's an array type
	if p.is(token.LEFT_BRACKET) {
		p.next()
		end := p.consumeExpected(token.RIGHT_BRACKET)

		return &ast.FlowArrayType{
			End:         end,
			ElementType: flowType,
		}
	}

	// it's a flow function type
	if !p.forbidUnparenthesizedFunctionType && p.is(token.ARROW) {
		p.next()

		returnType := p.parseFlowType()

		return &ast.FlowFunctionType{
			Start: start,
			Parameters: []*ast.FlowFunctionParameter{
				{Type: flowType},
			},
			ReturnType: returnType,
		}
	}

	if p.is(token.OR) && p.scope.allowUnionType {
		return p.parseFlowUnionType(flowType)
	}

	if p.is(token.AND) && p.scope.allowIntersectionType {
		return p.parseFlowIntersectionType(flowType)
	}

	if flowType == nil {
		p.next()
		return nil
	}

	return flowType
}

func (p *Parser) parseFlowTypeAnnotation() ast.FlowType {
	closeTypeScope := p.openTypeScope()
	defer closeTypeScope()

	p.consumeExpected(token.COLON)

	return p.parseFlowType()
}
