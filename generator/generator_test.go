package generator

import (
	"testing"
	"yawp/options"
	"yawp/parser"
)

func TestPlayground(t *testing.T) {
	opt := &options.Options{
		Target: options.ES2020,
		Minify: true,
	}

	// language=js
	prog, err := parser.ParseFile("", `
		const test = 132, foo = '';
		test.a = 123;
	`)

	if err != nil {
		t.Fail()
	}

	str := Generate(opt, prog)

	g := newIdGenerator()
	gs := make([]string, 1000)

	for i:=0;i<1000;i++ {
		gs[i] = g.Next()
	}

	_ = str
	return
}
