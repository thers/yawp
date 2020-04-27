package parser

import (
	"fmt"
	"path"
	"path/filepath"
	"runtime"
	"testing"
)

// Quick and dirty replacement for terst

func tt(t *testing.T, f func()) {
	defer func() {
		if x := recover(); x != nil {
			_, file, line, _ := runtime.Caller(4)
			t.Errorf("Error at %s:%d: %v", filepath.Base(file), line, x)
		}
	}()

	f()
}

func makeAssert(t *testing.T) func(string, interface{}) {
	return func(src string, expected interface{}) {
		_, err := ParseModule("", src)

		if err == nil && expected == nil {
			return
		}

		if expected == nil {
			if err == nil {
				return
			}
		}

		var actual interface{}

		if err == nil {
			actual = nil
		} else if errString, ok := err.(error); ok {
			actual = errString.Error()
		}

		if expected != actual {
			name := t.Name()
			_, filePath, line, _ := runtime.Caller(1)
			_, file := path.Split(filePath)

			fmt.Printf("\t%s: %s:%d: Expected output: %v\n", name, file, line, expected)
			fmt.Printf("\t%s: %s:%d: Actual output: %v\n", name, file, line, actual)

			t.Fail()
		}
	}
}

func is(a, b interface{}) {
	as := fmt.Sprintf("%v", a)
	bs := fmt.Sprintf("%v", b)
	if as != bs {
		panic(fmt.Errorf("%+v(%T) != %+v(%T)", a, a, b, b))
	}
}
