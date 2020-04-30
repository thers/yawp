package generator

import (
	"strings"
	"yawp/ids"
	"yawp/options"
	"yawp/parser/ast"
	"yawp/parser/file"
)

var ident = "                                                                                                          "

type Generator struct {
	*ast.Walker

	source  string
	options *options.Options

	output *strings.Builder
	ids    *ids.Ids

	identLevel int

	// punctuation flags
	wrapExpression bool
}

func (g *Generator) str(s string) *Generator {
	g.output.WriteString(s)

	return g
}

func (g *Generator) rune(r rune) *Generator {
	g.output.WriteRune(r)

	return g
}

func (g *Generator) semicolon() *Generator {
	g.output.WriteRune(';')

	return g
}

func (g *Generator) nl() *Generator {
	g.output.WriteRune('\n')
	g.ident()

	return g
}

func (g *Generator) src(loc *file.Loc) *Generator {
	g.output.WriteString(g.source[loc.From:loc.To])

	return g
}

func (g *Generator) ident() *Generator {
	g.output.WriteString(ident[:g.identLevel])

	return g
}

func (g *Generator) indentInc() *Generator {
	g.identLevel += 2

	if g.identLevel > len(ident) {
		ident += ident
	}

	return g
}

func (g *Generator) indentDec() *Generator {
	g.identLevel -= 2

	if g.identLevel < 0 {
		g.identLevel = 0
	}

	return g
}

func Generate(options *options.Options, program *ast.Module) string {
	generator := &Generator{
		Walker: &ast.Walker{},

		source:     program.File.Source(),
		output:     &strings.Builder{},
		options:    options,
		ids:        program.Ids,
		identLevel: 0,
	}
	generator.Walker.Visitor = generator

	program.Visit(generator)

	return generator.output.String()
}
