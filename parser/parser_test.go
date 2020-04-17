package parser

import (
	"testing"

	"yawp/parser/ast"
)

func testParse(src string) (parser *Parser, program *ast.Program, err error) {
	defer func() {
		if tmp := recover(); tmp != nil {
			panic(tmp)
		}
	}()
	parser = newParser("", src)
	program, err = parser.parse()
	return
}

func TestErrorsLocations(t *testing.T) {
	assert := makeAssert(t)

	assert("a", nil)
	assert("#a", "1:1 Unexpected token #")
	assert("\n\n\n#a", "4:1 Unexpected token #")
	assert("\n\n1+2+3\t+#", "3:8 Unexpected token #")
}

func TestErrors(t *testing.T) {
	assert := makeAssert(t)

	assert("!{ *a(b, b){} };", "1:4 Unexpected token *")
	assert("var x, ;", "1:8 Unexpected token ;")
	assert("function*g() { try {} catch (yield) {} }", "1:30 Unexpected token yield")
}

func TestJSXAmbiguities(t *testing.T) {
	assert := makeAssert(t)

	assert(`t=2<2`, nil)
	assert(`t=<a />`, nil)

	// one may may wonder "why is it even throw error?!"
	// well, it's not easy:
	// 1) Flow parser does exactly the same
	// 2) Babel can successfully parse it, but it uses rather complex process to disambiguate
	// it's a pretty rare case so we don't even try to parse it correctly in favor of speed
	assert(`t=<a>()=>null;t='</a>'`, "1:22 Unexpected token ILLEGAL")

	assert(`foo?.2:bar`, nil)
	assert(`foo?.bar`, nil)

	assert(`t = async()`, nil)
	assert(`t = async()=>null`, nil)
	assert(`t = async await=>null`, nil)

	assert(`t = (foo)`, nil)
	assert(`t = (foo) => bar`, nil)

	assert(`type T = | a & b | c`, nil)
	assert(`type T = a & (b | f) | c`, nil)
}
