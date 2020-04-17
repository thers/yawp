package parser

import "testing"

func TestPlayground(t *testing.T) {
	// language=js
	src := `
		t=<a>()=>null;t='</a>
`

	prog, err := ParseFile("", src)

	if err != nil {
		panic(err)
	}

	body := prog.Body
	_ = body
	return
}
