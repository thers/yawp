package parser

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"

	"yawp/parser/ast"
	"yawp/parser/file"
	"yawp/parser/token"
)

type Parser struct {
	parsedStr string
	src       string
	length    int
	base      int

	col         int   // column of the current char
	line        int   // line of the current char

	chr           rune // The current character
	chrOffset     int  // The offset of current character
	nextChrOffset int  // The nextChrOffset after current character (may be greater than 1)

	idx            file.Idx    // location of token
	token          token.Token // token itself
	literal        string      // literal of a token, if any
	tokenIsKeyword bool        // is current token a keyword

	scope *Scope

	insertSemicolon   bool // If we see a newline, then insert an implicit semicolon
	implicitSemicolon bool // An implicit semicolon exists

	genericTypeParametersMode bool

	forbidUnparenthesizedFunctionType bool

	jsxTextParseFrom int

	errors ErrorList

	file *file.File
}

func newParser(filename, src string, base int) *Parser {
	return &Parser{
		chr:    ' ', // This is set so we can start scanning by skipping whitespace
		col:    1,
		line:   1,
		src:    src,
		length: len(src),
		base:   base,
		file:   file.NewFile(filename, src, base),
	}
}

func ReadSource(filename string, src interface{}) ([]byte, error) {
	if src != nil {
		switch src := src.(type) {
		case string:
			return []byte(src), nil
		case []byte:
			return src, nil
		case *bytes.Buffer:
			if src != nil {
				return src.Bytes(), nil
			}
		case io.Reader:
			var bfr bytes.Buffer
			if _, err := io.Copy(&bfr, src); err != nil {
				return nil, err
			}
			return bfr.Bytes(), nil
		}
		return nil, errors.New("invalid source")
	}
	return ioutil.ReadFile(filename)
}

// ParseFile parses the source code of a single JavaScript/ECMAScript source file and returns
// the corresponding ast.Program node.
//
// If fileSet == nil, ParseFile parses source without a FileSet.
// If fileSet != nil, ParseFile first adds filename and src to fileSet.
//
// The filename argument is optional and is used for labelling errors, etc.
//
// src may be a string, a byte slice, a bytes.Buffer, or an io.Reader, but it MUST always be in UTF-8.
//
//      // Parse some JavaScript, yielding a *ast.Program and/or an ErrorList
//      program, err := parser.ParseFile(nil, "", `if (abc > 1) {}`, 0)
//
func ParseFile(fileSet *file.FileSet, filename string, src interface{}) (*ast.Program, error) {
	str, err := ReadSource(filename, src)
	if err != nil {
		return nil, err
	}
	{
		str := string(str)

		base := 1
		if fileSet != nil {
			base = fileSet.AddFile(filename, str)
		}

		parser := newParser(filename, str, base)
		prog, perr := parser.parse()

		return prog, perr
	}
}

func (p *Parser) parse() (*ast.Program, error) {
	p.next()
	program := p.parseProgram()

	return program, p.errors.Err()
}

func (p *Parser) next() (idx file.Idx) {
	idx = p.idx
	p.token, p.literal, p.idx = p.scan()
	return
}

// ParseFunctionForTests parses a given parameter list and body as a function and returns the
// corresponding ast.FunctionLiteral node.
//
// The parameter list, if any, should be a comma-separated list of identifiers.
//
func ParseFunctionForTests(parameterList, body string) (*ast.FunctionLiteral, error) {

	src := "(function(" + parameterList + ") {\n" + body + "\n})"

	parser := newParser("", src, 1)
	program, err := parser.parse()
	if err != nil {
		return nil, err
	}

	return program.Body[0].(*ast.ExpressionStatement).Expression.(*ast.FunctionLiteral), nil
}
