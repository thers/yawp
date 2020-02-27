package main

import (
	"fmt"
	"yawp/parser"
)

func main() {
	filename := "" // A filename is optional
	src := `
    // Sample xyzzy example
    (function(){
        if (3.14159 > 0) {
            console.log(..."Hello, World.");
            return;
        }

        var xyzzy = NaN;
        console.log("Nothing happens.");
        return xyzzy;
    })();
`

	// Parse some JavaScript, yielding a *ast.Program and/or an ErrorList
	program, err := parser.ParseFile(nil, filename, src, 0)
	if err != nil {
		panic(err)
	}

	fmt.Print(program)

	return

	//tokenizer, terr := tokenizer.NewTokenizer(os.Args[1])
	//if terr != nil {
	//	panic(terr)
	//}
	//
	//s := time.Now()
	//tokens, err := tokenizer.Parse()
	//if err != nil {
	//	panic(err)
	//}
	//
	//e := time.Since(s)
	//
	//tokensJson, _ := json.MarshalIndent(tokens, "", "  ")
	//fmt.Printf("%s\nTook: %s\n", tokensJson, e)
}
