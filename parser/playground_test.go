package parser

import "testing"

func TestPlayground(t *testing.T) {
	// language=js
	src := `
		while(true){true;false}
`

	prog, err := ParseFile("", src)

	if err != nil {
		panic(err)
	}

	body := prog.Body
	_ = body
	return
}
