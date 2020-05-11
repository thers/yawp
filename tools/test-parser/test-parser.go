package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"yawp/parser"
)

type runnerFn func(string) bool

/*
This will run official es262 parser test cases
and report if anything is wrong
*/
func main() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	rootPath := path.Join(cwd, "test262-parser-tests")
	passDone := make(chan bool)
	passesDone := 0

	go runFiles(path.Join(rootPath, "pass"), passDone, func(file string) bool {
		if shouldSkipPass(file) {
			return true
		}

		return assertFile(file, false)
	})

	go runFiles(path.Join(rootPath, "fail"), passDone, func(file string) bool {
		//if shouldSkipPass(file) {
		//	return
		//}

		return assertFile(file, true)
	})

	select {
	case <-passDone:
		passesDone++

		if passesDone == 2 {
			return
		}
	}
}

func runFiles(path string, done chan bool, runner runnerFn) {
	passFiles := readDir(path)

	totalCnt := len(passFiles)
	fileCnt := 0
	killed := false

	for _, passFile := range passFiles {
		fileCnt++

		func(file string, cnt int){
			if killed {
				return
			}

			if !runner(file) {
				killed = true
				done <- true
				return
			}

			if killed {
				return
			}

			if cnt == totalCnt {
				done <- true
			}
		}(passFile, fileCnt)
	}
}

func shouldSkipPass(file string) bool {
	fileId := path.Base(file)

	switch fileId {
	// 0008, test for parsing deprecated in non-strict-mode octal literal
	// since we assume everything a strict-mode code this test is pointless
	case "0b6dfcd5427a43a6.js", "45dd9586f26a3cf4.js", "7c6d13458e08e1f4.js", "84b2a5d834daee2f.js",
		"a953f09a1b6b6725.js", "bf6aaaab7c143ca1.js", "d2af344779cc1f26.js", "f1534392279bddbf.js":
		fallthrough
	// let is reserved keyword in strict-mode
	case "14199f22a45c7e30.js", "2ef5ba0343d739dc.js", "4f731d62a74ab666.js", "5654d4106d7025c2.js",
		"56e2ba90e05f5659.js", "5ecbbdc097bee212.js", "63c92209eb77315a.js", "65401ed8dc152370.js",
		"660f5a175a2d46ac.js", "6815ab22de966de8.js", "6b36b5ad4f3ad84d.js", "818ea8eaeef8b3da.js",
		"9aa93e1e417ce8e3.js", "9fe1d41db318afba.js", "a1594a4d0c0ee99a.js", "b8c98b5cd38f2bd9.js",
		"c442dc81201e2b55.js", "c8565124aee75c69.js", "df696c501125c86f.js", "ee4e8fa6257d810a.js",
		"f0d9a7a2f5d42210.js", "f0fbbdabdaca2146.js", "f2e41488e95243a8.js", "ffaf5b9d3140465b.js":
		fallthrough
	// interface, yield, static too
	case "a91ad31c88855e59.js", "08c3105bb3f7ccb7.js", "09e84f25af85b836.js", "0d137e8a97ffe083.js",
		"12556d5e39db1cea.js", "194b702816a7e5e5.js", "19d1d07fe88ec849.js", "213c3b05c6690d2d.js",
		"26b946d7cc01c226.js", "29e41f46ede71f11.js", "3bbd75d597d54fe6.js", "438521c40cf1b08b.js",
		"5856de37689f8db9.js", "5ed18bdbe48cc4c3.js", "6196b3f969486455.js", "6776e2c88e03f1a2.js",
		"6b76b8761a049c19.js", "75172741c27c7703.js", "753341b6f22ec11f.js", "7993945fc0f58feb.js",
		"8d2ffebc7214c34f.js", "8de47a8b53495792.js", "90abe26c46af6975.js", "979b36a2c530f286.js",
		"9a666205cafd530f.js", "a157424306915066.js", "a2cb5a14559c6a50.js", "a93f6b22796d4868.js",
		"af97a3752e579223.js", "b0c6752e1db068ed.js", "b0d44fd20353fd82.js", "b205355de22689d1.js",
		"c086a8a5c8ef2bb9.js", "ce5f3bc27d5ccaac.js", "d22f8660531e1c1a.js", "d82ae3dbc61808f8.js",
		"e7c1f6f0913c4a95.js", "ec99a663d6f3983d.js", "f4a61fcdefebb9d4.js", "fd167642d02f2c66.js",
		"ffcf0064736d41e7.js":
		fallthrough
	// html comments are not support
	case "8ec6a55806087669.js", "946bee37652a31fa.js", "9f0d8eb6f7ab8180.js", "ba00173ff473e7da.js",
		"c532e126a986c1d4.js", "e03ae54743348d7d.js", "e2c80df1960433a3.js", "1270d541e0fd6af8.js",
		"4ae32442eef8a4e0.js", "4f5419fe648c691b.js", "5a2a8e992fa4fe37.js", "5d5b9de6d9b95f3e.js",
		"b15ab152f8531a9f.js", "d3ac25ddc7ba9779.js", "fbcd793ec7c82779.js":
		return true
	}

	return false
}

func assertFile(file string, expectError bool) bool {
	fileContent, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	defer func() {
		err := recover()

		if err != nil {

			fmt.Printf("File %s is......NOT OK\n", path.Base(file))
			printFile(file, fileContent)
			panic(err)
		}
	}()

	_, err = parser.ParseModule(file, fileContent)

	if err == nil && expectError {
		fmt.Printf("\nWas expecting error, but got nothing.\n")
		printFile(file, fileContent)
		return false
	}

	if err != nil && !expectError {
		fmt.Printf("\nWas NOT expecting error: %s\n", err)
		printFile(file, fileContent)
		printFileError(fileContent, err)
		return false
	}

	return true
}

func readDir(subdir string) []string {
	var err error
	var dir []os.FileInfo

	dir, err = ioutil.ReadDir(subdir)
	if err != nil {
		panic(err)
	}

	files := make([]string, 0, len(dir))

	for _, file := range dir {
		if file.IsDir() {
			continue
		}

		files = append(files, path.Join(subdir, file.Name()))
	}

	return files
}

func printFile(file string, fileContent []byte) {
	fmt.Printf("Content:\n\n%s\n", fileContent)
}

func printFileError(fileContent []byte, err error) {
	if se, ok := err.(*parser.SyntaxError); ok {
		for i := 0; i < int(se.Loc.From); i++ {
			fmt.Print("~")
		}

		fmt.Print("^\n")
	}
}
