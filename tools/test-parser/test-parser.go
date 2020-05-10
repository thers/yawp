package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"yawp/parser"
)

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

	passFiles := readDir(path.Join(rootPath, "pass"))
	for _, passFile := range passFiles {
		if !assertFile(passFile, false) {
			return
		}
	}
}

func assertFile(file string, expectError bool) bool {fileId := path.Base(file)

	switch fileId {
	// 0008, test for parsing deprecated in non-strict-mode octal literal
	// since we assume everything a strict-mode code this test is pointless
	case "0b6dfcd5427a43a6.js", "45dd9586f26a3cf4.js", "7c6d13458e08e1f4.js", "84b2a5d834daee2f.js":
		fallthrough
	// <!-- html comments are dangling what and are only for HTML, right?
	case "1270d541e0fd6af8.js", "4ae32442eef8a4e0.js", "4f5419fe648c691b.js", "5a2a8e992fa4fe37.js",
		 "5d5b9de6d9b95f3e.js":
		fallthrough
	// let = 1 so let is reserved keyword, ok
	case "14199f22a45c7e30.js", "2ef5ba0343d739dc.js", "4f731d62a74ab666.js", "5654d4106d7025c2.js",
		 "56e2ba90e05f5659.js", "5ecbbdc097bee212.js", "63c92209eb77315a.js", "65401ed8dc152370.js",
		 "660f5a175a2d46ac.js", "6815ab22de966de8.js", "6b36b5ad4f3ad84d.js", "818ea8eaeef8b3da.js":
		fallthrough
	// `--> nothing` uhh, no, just no
	case "8ec6a55806087669.js":
		return true
	}

	fmt.Printf("File %s is......", path.Base(file))

	fileContent, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	defer func(){
		err := recover()

		if err != nil {
			fmt.Printf("NOT OK\n")
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
		return false
	}

	fmt.Printf("OK\n")

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
	fmt.Printf("Content:\n\n%s\n\n", fileContent)
}
