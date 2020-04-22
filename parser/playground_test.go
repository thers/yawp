package parser

import "testing"

func TestPlayground(t *testing.T) {
	// language=js
	src := `
		const test = 123;
		test = 321;
		
		if (true) {
		    const test = 321;
		    test++;
		}
`

	prog, err := ParseFile("", src)

	if err != nil {
		panic(err)
	}

	body := prog.Body
	_ = body
	return
}
