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
		Target: options.ES5,
		Minify: false,
	}

	parserStart := time.Now()
	// language=js
	prog, err := parser.ParseModule("", `
a = b-->1;
 --> nothing
	`)
	fmt.Println("Parser pass took:", time.Since(parserStart))

	if err != nil || prog == nil {
		t.Fail()
	}

	optimizerStart := time.Now()
	optimizer.Optimize(prog, opt)
	fmt.Println("Optimizer pass took: ", time.Since(optimizerStart))

	generatorStart := time.Now()
	str := Generate(opt, prog)
	fmt.Println("Generator pass took:", time.Since(generatorStart))

	_ = str
	fmt.Println(str)
	return
}
