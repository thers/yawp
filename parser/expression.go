package parser

import (
	"yawp/parser/ast"
	"yawp/parser/file"
	"yawp/parser/token"
)

func (p *Parser) parseIdentifier() *ast.Identifier {
	literal := p.literal
	idx := p.idx
	p.next()
	return &ast.Identifier{
		Name: literal,
		Idx:  idx,
	}
}

func (p *Parser) parsePrimaryExpression() ast.Expression {
	literal := p.literal
	idx := p.idx

	switch p.token {
	case token.IDENTIFIER:
		return p.parseIdentifierOrSingleArgumentArrowFunction()
	case token.NULL:
		p.next()
		return &ast.NullLiteral{
			Idx:     idx,
			Literal: literal,
		}
	case token.BOOLEAN:
		p.next()
		value := false
		switch literal {
		case "true":
			value = true
		case "false":
			value = false
		default:
			p.error(idx, "Illegal boolean literal")
		}
		return &ast.BooleanLiteral{
			Idx:     idx,
			Literal: literal,
			Value:   value,
		}
	case token.STRING:
		p.next()
		value, err := parseStringLiteral(literal[1 : len(literal)-1])
		if err != nil {
			p.error(idx, err.Error())
		}
		return &ast.StringLiteral{
			Idx:     idx,
			Literal: literal,
			Value:   value,
		}
	case token.NUMBER:
		p.next()
		value, err := parseNumberLiteral(literal)
		if err != nil {
			p.error(idx, err.Error())
			value = 0
		}
		return &ast.NumberLiteral{
			Idx:     idx,
			Literal: literal,
			Value:   value,
		}
	case token.SLASH, token.QUOTIENT_ASSIGN:
		return p.parseRegExpLiteral()
	case token.LEFT_BRACE:
		return p.parseObjectLiteral()
	case token.LEFT_BRACKET:
		return p.parseArrayLiteral()
	case token.LEFT_PARENTHESIS:
		return p.parseArrowFunctionOrSequenceExpression()
	case token.THIS:
		p.next()
		return &ast.ThisExpression{
			Idx: idx,
		}
	case token.FUNCTION:
		return p.parseFunction(false)
	}

	p.errorUnexpectedToken(p.token)
	p.nextStatement()
	return &ast.BadExpression{From: idx, To: p.idx}
}

func (p *Parser) parseRegExpLiteral() *ast.RegExpLiteral {

	offset := p.chrOffset - 1 // Opening slash already gotten
	if p.is(token.QUOTIENT_ASSIGN) {
		offset -= 1 // =
	}
	idx := p.idxOf(offset)

	pattern, err := p.scanString(offset)
	endOffset := p.chrOffset

	if err == nil {
		pattern = pattern[1 : len(pattern)-1]
	}

	flags := ""
	if !isLineTerminator(p.nextChr) && !isLineWhiteSpace(p.nextChr) {
		p.next()

		if p.is(token.IDENTIFIER) { // gim

			flags = p.literal
			p.next()
			endOffset = p.chrOffset - 1
		}
	} else {
		p.next()
	}

	literal := p.str[offset:endOffset]

	return &ast.RegExpLiteral{
		Idx:     idx,
		Literal: literal,
		Pattern: pattern,
		Flags:   flags,
	}
}

func (p *Parser) parseArrayLiteral() ast.Expression {

	idx0 := p.consumeExpected(token.LEFT_BRACKET)
	var value []ast.Expression
	for !p.is(token.RIGHT_BRACKET) && !p.is(token.EOF) {
		if p.is(token.COMMA) {
			p.next()
			value = append(value, nil)
			continue
		}
		value = append(value, p.parseAssignmentExpression())
		if !p.is(token.RIGHT_BRACKET) {
			p.consumeExpected(token.COMMA)
		}
	}
	idx1 := p.consumeExpected(token.RIGHT_BRACKET)

	return &ast.ArrayLiteral{
		LeftBracket:  idx0,
		RightBracket: idx1,
		Value:        value,
	}
}

func (p *Parser) parseArgumentList() (argumentList []ast.Expression, idx0, idx1 file.Idx) {
	idx0 = p.consumeExpected(token.LEFT_PARENTHESIS)
	if !p.is(token.RIGHT_PARENTHESIS) {
		for {
			argumentList = append(argumentList, p.parseAssignmentExpression())
			if !p.is(token.COMMA) {
				break
			}
			p.next()
		}
	}
	idx1 = p.consumeExpected(token.RIGHT_PARENTHESIS)
	return
}

func (p *Parser) parseCallExpression(left ast.Expression) ast.Expression {
	argumentList, idx0, idx1 := p.parseArgumentList()
	return &ast.CallExpression{
		Callee:           left,
		LeftParenthesis:  idx0,
		ArgumentList:     argumentList,
		RightParenthesis: idx1,
	}
}

func (p *Parser) parseDotMember(left ast.Expression) ast.Expression {
	period := p.consumeExpected(token.PERIOD)

	literal := p.literal
	idx := p.idx

	if !matchIdentifier.MatchString(literal) {
		p.consumeExpected(token.IDENTIFIER)
		p.nextStatement()
		return &ast.BadExpression{From: period, To: p.idx}
	}

	p.next()

	return &ast.DotExpression{
		Left: left,
		Identifier: ast.Identifier{
			Idx:  idx,
			Name: literal,
		},
	}
}

func (p *Parser) parseBracketMember(left ast.Expression) ast.Expression {
	idx0 := p.consumeExpected(token.LEFT_BRACKET)
	member := p.parseExpression()
	idx1 := p.consumeExpected(token.RIGHT_BRACKET)
	return &ast.BracketExpression{
		LeftBracket:  idx0,
		Left:         left,
		Member:       member,
		RightBracket: idx1,
	}
}

func (p *Parser) parseNewExpression() ast.Expression {
	idx := p.consumeExpected(token.NEW)
	callee := p.parseLeftHandSideExpression()
	node := &ast.NewExpression{
		New:    idx,
		Callee: callee,
	}
	if p.is(token.LEFT_PARENTHESIS) {
		argumentList, idx0, idx1 := p.parseArgumentList()
		node.ArgumentList = argumentList
		node.LeftParenthesis = idx0
		node.RightParenthesis = idx1
	}
	return node
}

func (p *Parser) parseLeftHandSideExpression() ast.Expression {

	var left ast.Expression
	if p.is(token.NEW) {
		left = p.parseNewExpression()
	} else {
		left = p.parsePrimaryExpression()
	}

	for {
		if p.is(token.PERIOD) {
			left = p.parseDotMember(left)
		} else if p.is(token.LEFT_BRACKET) {
			left = p.parseBracketMember(left)
		} else {
			break
		}
	}

	return left
}

func (p *Parser) parseLeftHandSideExpressionAllowCall() ast.Expression {

	allowIn := p.scope.allowIn
	p.scope.allowIn = true
	defer func() {
		p.scope.allowIn = allowIn
	}()

	var left ast.Expression
	if p.is(token.NEW) {
		left = p.parseNewExpression()
	} else {
		left = p.parsePrimaryExpression()
	}

	for {
		if p.is(token.PERIOD) {
			left = p.parseDotMember(left)
		} else if p.is(token.LEFT_BRACKET) {
			left = p.parseBracketMember(left)
		} else if p.is(token.LEFT_PARENTHESIS) {
			left = p.parseCallExpression(left)
		} else {
			break
		}
	}

	return left
}

func (p *Parser) parsePostfixExpression() ast.Expression {
	operand := p.parseLeftHandSideExpressionAllowCall()

	switch p.token {
	case token.INCREMENT, token.DECREMENT:
		// Make sure there is no line terminator here
		if p.implicitSemicolon {
			break
		}
		tkn := p.token
		idx := p.idx
		p.next()
		switch operand.(type) {
		case *ast.Identifier, *ast.DotExpression, *ast.BracketExpression:
		default:
			p.error(idx, "Invalid left-hand side in assignment")
			p.nextStatement()
			return &ast.BadExpression{From: idx, To: p.idx}
		}
		return &ast.UnaryExpression{
			Operator: tkn,
			Idx:      idx,
			Operand:  operand,
			Postfix:  true,
		}
	}

	return operand
}

func (p *Parser) parseUnaryExpression() ast.Expression {

	switch p.token {
	case token.PLUS, token.MINUS, token.NOT, token.BITWISE_NOT:
		fallthrough
	case token.DELETE, token.VOID, token.TYPEOF:
		tkn := p.token
		idx := p.idx
		p.next()
		return &ast.UnaryExpression{
			Operator: tkn,
			Idx:      idx,
			Operand:  p.parseUnaryExpression(),
		}
	case token.INCREMENT, token.DECREMENT:
		tkn := p.token
		idx := p.idx
		p.next()
		operand := p.parseUnaryExpression()
		switch operand.(type) {
		case *ast.Identifier, *ast.DotExpression, *ast.BracketExpression:
		default:
			p.error(idx, "Invalid left-hand side in assignment")
			p.nextStatement()
			return &ast.BadExpression{From: idx, To: p.idx}
		}
		return &ast.UnaryExpression{
			Operator: tkn,
			Idx:      idx,
			Operand:  operand,
		}
	}

	return p.parsePostfixExpression()
}

func (p *Parser) parseMultiplicativeExpression() ast.Expression {
	next := p.parseUnaryExpression
	left := next()

	for p.is(token.MULTIPLY) || p.is(token.SLASH) ||
		p.is(token.REMAINDER) {
		tkn := p.token
		p.next()
		left = &ast.BinaryExpression{
			Operator: tkn,
			Left:     left,
			Right:    next(),
		}
	}

	return left
}

func (p *Parser) parseAdditiveExpression() ast.Expression {
	next := p.parseMultiplicativeExpression
	left := next()

	for p.is(token.PLUS) || p.is(token.MINUS) {
		tkn := p.token
		p.next()
		left = &ast.BinaryExpression{
			Operator: tkn,
			Left:     left,
			Right:    next(),
		}
	}

	return left
}

func (p *Parser) parseShiftExpression() ast.Expression {
	next := p.parseAdditiveExpression
	left := next()

	for p.is(token.SHIFT_LEFT) || p.is(token.SHIFT_RIGHT) ||
		p.is(token.UNSIGNED_SHIFT_RIGHT) {
		tkn := p.token
		p.next()
		left = &ast.BinaryExpression{
			Operator: tkn,
			Left:     left,
			Right:    next(),
		}
	}

	return left
}

func (p *Parser) parseRelationalExpression() ast.Expression {
	next := p.parseShiftExpression
	left := next()

	allowIn := p.scope.allowIn
	p.scope.allowIn = true
	defer func() {
		p.scope.allowIn = allowIn
	}()

	switch p.token {
	case token.LESS, token.LESS_OR_EQUAL, token.GREATER, token.GREATER_OR_EQUAL:
		tkn := p.token
		p.next()
		return &ast.BinaryExpression{
			Operator:   tkn,
			Left:       left,
			Right:      p.parseRelationalExpression(),
			Comparison: true,
		}
	case token.INSTANCEOF:
		tkn := p.token
		p.next()
		return &ast.BinaryExpression{
			Operator: tkn,
			Left:     left,
			Right:    p.parseRelationalExpression(),
		}
	case token.IN:
		if !allowIn {
			return left
		}
		tkn := p.token
		p.next()
		return &ast.BinaryExpression{
			Operator: tkn,
			Left:     left,
			Right:    p.parseRelationalExpression(),
		}
	}

	return left
}

func (p *Parser) parseEqualityExpression() ast.Expression {
	next := p.parseRelationalExpression
	left := next()

	for p.is(token.EQUAL) || p.is(token.NOT_EQUAL) ||
		p.is(token.STRICT_EQUAL) || p.is(token.STRICT_NOT_EQUAL) {
		tkn := p.token
		p.next()
		left = &ast.BinaryExpression{
			Operator:   tkn,
			Left:       left,
			Right:      next(),
			Comparison: true,
		}
	}

	return left
}

func (p *Parser) parseBitwiseAndExpression() ast.Expression {
	next := p.parseEqualityExpression
	left := next()

	for p.is(token.AND) {
		tkn := p.token
		p.next()
		left = &ast.BinaryExpression{
			Operator: tkn,
			Left:     left,
			Right:    next(),
		}
	}

	return left
}

func (p *Parser) parseBitwiseExclusiveOrExpression() ast.Expression {
	next := p.parseBitwiseAndExpression
	left := next()

	for p.is(token.EXCLUSIVE_OR) {
		tkn := p.token
		p.next()
		left = &ast.BinaryExpression{
			Operator: tkn,
			Left:     left,
			Right:    next(),
		}
	}

	return left
}

func (p *Parser) parseBitwiseOrExpression() ast.Expression {
	next := p.parseBitwiseExclusiveOrExpression
	left := next()

	for p.is(token.OR) {
		tkn := p.token
		p.next()
		left = &ast.BinaryExpression{
			Operator: tkn,
			Left:     left,
			Right:    next(),
		}
	}

	return left
}

func (p *Parser) parseLogicalAndExpression() ast.Expression {
	next := p.parseBitwiseOrExpression
	left := next()

	for p.is(token.LOGICAL_AND) {
		tkn := p.token
		p.next()
		left = &ast.BinaryExpression{
			Operator: tkn,
			Left:     left,
			Right:    next(),
		}
	}

	return left
}

func (p *Parser) parseLogicalOrExpression() ast.Expression {
	next := p.parseLogicalAndExpression
	left := next()

	for p.is(token.LOGICAL_OR) {
		tkn := p.token
		p.next()
		left = &ast.BinaryExpression{
			Operator: tkn,
			Left:     left,
			Right:    next(),
		}
	}

	return left
}

func (p *Parser) parseConditionalExpression() ast.Expression {
	left := p.parseLogicalOrExpression()

	if p.is(token.QUESTION_MARK) {
		p.next()
		consequent := p.parseAssignmentExpression()
		p.consumeExpected(token.COLON)
		return &ast.ConditionalExpression{
			Test:       left,
			Consequent: consequent,
			Alternate:  p.parseAssignmentExpression(),
		}
	}

	return left
}

func (p *Parser) parseAssignmentExpression() ast.Expression {
	var operator token.Token

	left := p.parseConditionalExpression()

	switch p.token {
	case token.ASSIGN:
		operator = p.token
	case token.ADD_ASSIGN:
		operator = token.PLUS
	case token.SUBTRACT_ASSIGN:
		operator = token.MINUS
	case token.MULTIPLY_ASSIGN:
		operator = token.MULTIPLY
	case token.QUOTIENT_ASSIGN:
		operator = token.SLASH
	case token.REMAINDER_ASSIGN:
		operator = token.REMAINDER
	case token.AND_ASSIGN:
		operator = token.AND
	case token.AND_NOT_ASSIGN:
		operator = token.AND_NOT
	case token.OR_ASSIGN:
		operator = token.OR
	case token.EXCLUSIVE_OR_ASSIGN:
		operator = token.EXCLUSIVE_OR
	case token.SHIFT_LEFT_ASSIGN:
		operator = token.SHIFT_LEFT
	case token.SHIFT_RIGHT_ASSIGN:
		operator = token.SHIFT_RIGHT
	case token.UNSIGNED_SHIFT_RIGHT_ASSIGN:
		operator = token.UNSIGNED_SHIFT_RIGHT
	}

	if operator != 0 {
		idx := p.idx
		p.next()
		switch left.(type) {
		case *ast.Identifier, *ast.DotExpression, *ast.BracketExpression:
		default:
			p.error(left.Idx0(), "Invalid left-hand side in assignment")
			p.nextStatement()
			return &ast.BadExpression{From: idx, To: p.idx}
		}
		return &ast.AssignExpression{
			Left:     left,
			Operator: operator,
			Right:    p.parseAssignmentExpression(),
		}
	}

	return left
}

func (p *Parser) parseExpression() ast.Expression {
	next := p.parseAssignmentExpression
	left := next()

	if p.is(token.COMMA) {
		sequence := []ast.Expression{left}
		for {
			if !p.is(token.COMMA) {
				break
			}
			p.next()
			sequence = append(sequence, next())
		}
		return &ast.SequenceExpression{
			Sequence: sequence,
		}
	}

	return left
}
