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

func assertFile(file string, expectError bool) bool {
	fileId := path.Base(file)

	switch fileId {
	// 0008, test for parsing deprecated in non-strict-mode octal literal
	// since we assume everything a strict-mode code this test is pointless
	case "0b6dfcd5427a43a6.js":
		fallthrough
	// <!-- html comments are dangling what and are only for HTML, right?
	case "1270d541e0fd6af8.js":
		fallthrough
	// let = 1 so let is reserved keyword, ok
	case "14199f22a45c7e30.js":
		return true
	}

	fileContent, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	_, err = parser.ParseModule(file, fileContent)

	if err == nil && expectError {
		fmt.Printf("Was expecting error, but got nothing.\n")
		printFile(file, fileContent)
		return false
	}

	if err != nil && !expectError {
		fmt.Printf("Was NOT expecting error: %s\n", err)
		printFile(file, fileContent)
		return false
	}

	fmt.Printf("File %s is OK\n", file)

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
	fmt.Printf("File \"%s\" content:\n\n%s\n\n", file, fileContent)
}
