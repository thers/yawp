package parser

import "yawp/parser/ast"

func (p *Parser) useSymbolsScope(stype ast.SymbolsScopeType) {
	symbolsScope := &ast.SymbolsScope{
		Type:     stype,
		Symbols:  make([]*ast.Symbol, 0),
		Parent:   p.symbolsScope,
		Children: make([]*ast.SymbolsScope, 0),
	}

	p.symbolsScope = symbolsScope
}

func (p *Parser) restoreSymbolsScope() {
	symbolsScope := p.symbolsScope.Parent

	if symbolsScope == nil {
		panic("Can not restore symbols scope as there's no parent scope")
	}

	symbolsScope.Children = append(symbolsScope.Children, p.symbolsScope)

	p.symbolsScope = symbolsScope
}

func (p *Parser) dropSymbolsScope() {
	symbolsScope := p.symbolsScope.Parent

	if symbolsScope == nil {
		panic("Can not drop symbols scope as there's no parent scope")
	}

	p.symbolsScope = symbolsScope
}


func (p *Parser) useSymbolFlags(flags ast.Flags) func() {
	prevFlags := p.symbolFlags
	p.symbolFlags = flags

	return func() {
		p.symbolFlags = prevFlags
	}
}

func (p *Parser) symbol(id *ast.Identifier, flags ast.Flags, stype ast.SymbolRefType) *ast.Identifier {
	id.Symbol = p.symbolsScope.MakeSymbol(id.Name)

	id.Symbol.Type = stype
	id.Symbol.Flags = flags

	return id
}