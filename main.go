package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime/pprof"
	"time"
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

	src, _ := ioutil.ReadFile(filename)

	f, err := os.Create("main.perf")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	start := time.Now()
	ast, err := parser.ParseFile(nil, filename, src)
	if err != nil {
		panic(err)
	}

	timeTaken := time.Since(start)
	throughput := float64((len(src)/1024.0)/1024.0)/timeTaken.Seconds()

	fmt.Printf("time taken: %s\n", time.Since(start))
	fmt.Printf("throughput: %fMB/s\n", throughput)
	fmt.Println(len(ast.Body))

	return
}
