package generator

import (
	"strings"
	"yawp/ids"
	"yawp/options"
	"yawp/parser/ast"
	"yawp/parser/file"
)

type Generator struct {
	*ast.DefaultVisitor

	source  string
	options *options.Options

	output *strings.Builder
	ids    *ids.Ids
}

func (g *Generator) str(s string) {
	g.output.WriteString(s)
}

func (g *Generator) rune(r rune) {
	g.output.WriteRune(r)
}

func (g *Generator) semicolon() {
	g.output.WriteRune(';')
}

func (g *Generator) nl() {
	g.output.WriteRune('\n')
}

func (g *Generator) src(loc *file.Loc) {
	g.output.WriteString(g.source[loc.From:loc.To])
}

func Generate(options *options.Options, program *ast.Module) string {
	generator := &Generator{
		DefaultVisitor: &ast.DefaultVisitor{},

		source:  program.File.Source(),
		output:  &strings.Builder{},
		options: options,
		ids:     program.Ids,
	}
	generator.DefaultVisitor.Specific = generator

	program.Visit(generator)

	return generator.output.String()
}
