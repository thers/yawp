package generator

import (
	"fmt"
	"testing"
	"time"
	"yawp/optimizer"
	"yawp/options"
	"yawp/parser"
)

func TestPlayground(t *testing.T) {
	opt := &options.Options{
		Target: options.ES2015,
		Minify: true,
	}

	// language=js
	prog, err := parser.ParseFile("", `
		const { t }= { t: 1 }
	`)

	if err != nil || prog == nil {
		t.Fail()
	}

	ops := time.Now()
	o := optimizer.NewOptimizer(prog, opt)
	prog.Visit(o)
	fmt.Println("Optimizer pass took: ", time.Since(ops))

	str := Generate(opt, prog)

	_ = str
	return
}
