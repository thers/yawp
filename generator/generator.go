package generator

import (
	"strings"
	"yawp/options"
	"yawp/parser/ast"
)

type Generator struct {
	opt *options.Options

	output *strings.Builder

	refScope *RefScope
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

func Generate(opt *options.Options, prog *ast.Program) string {
	generator := &Generator{
		output: &strings.Builder{},
		opt:    opt,
	}

	generator.program(prog)

	return generator.output.String()
}

func (g *Generator) program(program *ast.Program) {
	g.pushRefScope()
	g.statements(program.Body)
}
