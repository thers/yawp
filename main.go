package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
	"yawp/optimizer"
	"yawp/options"
	"yawp/parser"
)

//var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	var filename string

	if len(os.Args) > 1 {
		filename = os.Args[1]
	} else {
		filename = "test/short.js"
	}

	opt := &options.Options{
		Target: options.ES2020,
		Minify: true,
	}

	src, _ := ioutil.ReadFile(filename)

	//f, err := os.Create("main.perf")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//pprof.StartCPUProfile(f)
	//defer pprof.StopCPUProfile()

	start := time.Now()
	ast, err := parser.ParseFile(filename, src)
	if err != nil {
		panic(err)
	}

	ops := time.Now()
	o := optimizer.NewOptimizer(ast, opt)
	ast.Visit(o)
	fmt.Println("Optimizer pass took: ", time.Since(ops))

	fmt.Printf("Parser pass took: %s\n", time.Since(start))
	fmt.Println(len(ast.Body))

	return
}
