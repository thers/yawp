package ast

import (
	"yawp/options"
	"yawp/parser/token"
)

type SymbolsScopeType int
type SymbolRefType int

const (
	SSTModule SymbolsScopeType = iota
	SSTFunction
	SSTClass
	SSTBlock
)

const (
	SDeclaration Flags = 1 << iota
	SWrite
	SRead
)

const (
	SRUnknown SymbolRefType = iota
	SRVar
	SRLet
	SRConst
	SRClass
	SRImport
	SRExport
	SRFn
	SRFnParam
	SRLabel
	SRBuiltin
)

type SymbolRef struct {
	Name string
	Type SymbolRefType

	Usages int
	Reads  int
	Writes int

	ShadowsRef    *SymbolRef
	ShadowedByRef *SymbolRef

	Mangled bool
}

type Symbol struct {
	Name    string
	RefType SymbolRefType
	Ref     *SymbolRef
	Flags   Flags
}

type SymbolsScope struct {
	Opts *options.Options
	Type SymbolsScopeType

	Symbols []*Symbol
	Refs    map[string]*SymbolRef

	Parent   *SymbolsScope
	Children []*SymbolsScope
}

type Identifier struct {
	ExprNode
	Name      string
	LegacyRef *SymbolRef
	Symbol    *Symbol
}

func (s *SymbolsScope) AllocateSymbol(name string) *Symbol {
	symbol := &Symbol{
		Name:    name,
		RefType: 0,
		Ref:     nil,
		Flags:   0,
	}

	s.Symbols = append(s.Symbols, symbol)

	return symbol
}

func (s *SymbolsScope) getRef(name string) (ref *SymbolRef, nested bool) {
	var ok bool

	nested = false

	if ref, ok = s.Refs[name]; ok {
		return
	}

	if s.Parent != nil {
		ref, _ = s.Parent.getRef(name)
		nested = true
	}

	return
}

func (s *SymbolsScope) allocateRef(symbol *Symbol) *SymbolRef {
	ref := &SymbolRef{
		Name: symbol.Name,
		Type: symbol.RefType,
	}

	s.Refs[symbol.Name] = ref

	return ref
}

func (s *SymbolsScope) ReferenceSymbols() {
	s.Refs = make(map[string]*SymbolRef)

	// First reference current scope symbols
	for _, symbol := range s.Symbols {
		ref, fromParentScope := s.getRef(symbol.Name)

		if symbol.Flags.Has(SDeclaration) {
			// For declaration symbols we have to work with ref from this scope,
			// and not from parent
			// also ref could also be created for symbol.Name
			// due to the use-before-defined of functions and classes

			if ref != nil {
				if fromParentScope {
					// We're shadowing parent scope's ref
					newRef := s.allocateRef(symbol)
					newRef.ShadowsRef = ref
					ref.ShadowedByRef = newRef
					ref = newRef
				} else {
					ref.Type = symbol.RefType
				}
			} else {
				ref = s.allocateRef(symbol)
			}
		} else {
			if ref == nil {
				ref = s.allocateRef(symbol)
			} else {
				if ref.Type == SRUnknown && symbol.RefType != SRUnknown {
					ref.Type = symbol.RefType
				}
			}

			if symbol.Flags.Has(SRead) {
				ref.Reads++
			}

			if symbol.Flags.Has(SWrite) {
				ref.Writes++
			}

			ref.Usages++
		}

		symbol.Ref = ref
	}

	// Now ask children to do the same
	for _, childScope := range s.Children {
		childScope.ReferenceSymbols()
	}
}

func SymbolRefTypeFromToken(value token.Token) SymbolRefType {
	switch value {
	case token.VAR:
		return SRVar
	case token.CONST:
		return SRConst
	case token.LET:
		return SRLet
	default:
		return SRUnknown
	}
}
