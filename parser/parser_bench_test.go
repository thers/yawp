package parser

import (
	"testing"
	"yawp/parser/ast"
)

func testParse(src string) (parser *Parser, program *ast.Module, err error) {
	defer func() {
		if tmp := recover(); tmp != nil {
			panic(tmp)
		}
	}()
	parser = newParser("", src)
	program, err = parser.parse()
	return
}

func BenchmarkParser(t *testing.B) {
	for i := 0; i < t.N; i++ {
		testParse("()=><div/>")
	}
}