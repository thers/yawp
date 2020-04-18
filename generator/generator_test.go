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
		class Test<T> extends React.PureComponent<Props, OwnProps> {
    		foo: T
		}
	`)

	if err != nil {
		t.Fail()
	}

	str := Generate(opt, prog)

	_ = str
	return
}
