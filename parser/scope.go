package parser

import (
	"yawp/parser/ast"
)

type Scope struct {
	outer           *Scope
	allowIn         bool
	inClass			bool
	inIteration     bool
	inSwitch        bool
	inFunction      bool
	declarationList []ast.Declaration

	labels []string
}

func (p *Parser) openScope() {
	p.scope = &Scope{
		outer:   p.scope,
		allowIn: true,
	}
}

func (p *Parser) closeScope() {
	p.scope = p.scope.outer
}


func (p *Parser) openClassScope() func() {
	p.openScope()
	wasInClass := p.scope.inClass
	p.scope.inClass = true

	return func() {
		p.scope.inClass = wasInClass
		p.closeScope()
	}
}

func (p *Parser) openFunctionScope() func() {
	p.openScope()
	wasInFunction := p.scope.inFunction
	p.scope.inFunction = true

	return func() {
		p.scope.inFunction = wasInFunction
		p.closeScope()
	}
}


func (self *Scope) declare(declaration ast.Declaration) {
	self.declarationList = append(self.declarationList, declaration)
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
