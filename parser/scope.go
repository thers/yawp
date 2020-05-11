package parser

type Scope struct {
	outer *Scope

	allowIn            bool
	allowAwait         bool
	allowYield         bool
	allowTypeAssertion bool

	inClass     bool
	inIteration bool
	inSwitch    bool
	inFunction  bool
	inType      bool

	allowUnionType        bool
	allowIntersectionType bool

	labels []string
}

func (p *Parser) openScope() {
	p.scope = &Scope{
		outer:                 p.scope,
		allowIn:               true,
		allowTypeAssertion:    false,
		allowUnionType:        true,
		allowIntersectionType: true,
	}
}

func (p *Parser) closeScope() {
	p.scope = p.scope.outer
}

func (p *Parser) openClassScope() func() {
	wasAllowYield := p.scope.allowYield
	p.openScope()
	wasInClass := p.scope.inClass
	p.scope.inClass = true
	p.scope.allowYield = wasAllowYield

	return func() {
		p.scope.inClass = wasInClass
		p.closeScope()
	}
}

func (p *Parser) openFunctionScope(generator bool, async bool) func() {
	p.openScope()

	wasInFunction := p.scope.inFunction
	wasAllowAwait := p.scope.allowAwait
	wasAllowYield := p.scope.allowYield

	p.scope.inFunction = true
	p.scope.allowAwait = async
	p.scope.allowYield = generator

	return func() {
		p.scope.allowAwait = wasAllowAwait
		p.scope.inFunction = wasInFunction
		p.scope.allowYield = wasAllowYield
		p.closeScope()
	}
}

func (p *Parser) openTypeScope() func() {
	p.openScope()
	p.scope.inType = true

	return func() {
		p.scope.inType = false
		p.closeScope()
	}
}

func (self *Scope) hasLabel(name string) bool {
	for _, label := range self.labels {
		if label == name {
			return true
		}
	}
	if self.outer != nil && !self.inFunction {
		// Crossing a function boundary to look for a label is verboten
		return self.outer.hasLabel(name)
	}
	return false
}

func (s *Scope) inModuleRoot() bool {
	if s == nil {
		return true
	}

	return !s.inFunction && !s.inType && !s.inClass && !s.inIteration && !s.inSwitch
}
