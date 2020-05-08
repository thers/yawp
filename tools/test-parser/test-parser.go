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
	fileContent, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	_, err = parser.ParseModule(file, fileContent)

	if err == nil && expectError {
		fmt.Printf("Was expecting error, but got nothing.\n")
		printFile(fileContent)
		return false
	}

	if err != nil && !expectError {
		fmt.Printf("Was NOT expecting error: %s\n", err)
		printFile(fileContent)
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

func printFile(fileContent []byte) {
	fmt.Printf("File content:\n\n%s\n\n", fileContent)
}
