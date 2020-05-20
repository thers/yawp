package ast

type SymbolsScopeType int

const (
	SSTModule SymbolsScopeType = iota
	SSTFunction
	SSTClass
	SSTBlock
)

const (
	SymbolDeclaration Flags = 1 << iota
	SymbolWrite
	SymbolRead
)

type SymbolsScope struct {
	Type     SymbolsScopeType
	Symbols  []*Symbol
	Parent   *SymbolsScope
	Children []*SymbolsScope
}

type Symbol struct {
	Name  string
	Type  SymbolRefType
	Ref   *SymbolRef
	Flags Flags
}

type Identifier struct {
	ExprNode
	Name      string
	LegacyRef *SymbolRef
	Symbol    *Symbol
}

func (s *SymbolsScope) MakeSymbol(name string) *Symbol {
	symbol := &Symbol{
		Name:  name,
		Type:  0,
		Ref:   nil,
		Flags: 0,
	}

	s.Symbols = append(s.Symbols, symbol)

	return symbol
}
