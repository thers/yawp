package generator

import (
	"fmt"
	"testing"
	"time"
	"yawp/transpiler"
	"yawp/options"
	"yawp/parser"
)

func TestPlayground(t *testing.T) {
	opt := &options.Options{
		Target: options.ES5,
		Minify: true,
	}

	parserStart := time.Now()
	// language=js
	prog, err := parser.ParseModule("", `
function test(...a) {return a}
	`)
	fmt.Println("Parser pass took:", time.Since(parserStart))

	if err != nil || prog == nil {
		t.Fail()
	}

	optimizerStart := time.Now()
	transpiler.Transpile(prog, opt)
	fmt.Println("Transpiler pass took: ", time.Since(optimizerStart))

	generatorStart := time.Now()
	str := Generate(opt, prog)
	fmt.Println("Generator pass took:", time.Since(generatorStart))

	_ = str
	fmt.Println(str)
	return
}
