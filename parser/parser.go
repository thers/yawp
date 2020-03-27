/*
Package parser implements a parser for JavaScript.

    import (
        "yawp/parser/parser"
    )

Parse and return an AST

    filename := "" // A filename is optional
    src := `
        // Sample xyzzy example
        (function(){
            if (3.14159 > 0) {
                console.log("Hello, World.");
                return;
            }

            var xyzzy = NaN;
            console.log("Nothing happens.");
            return xyzzy;
        })();
    `

    // Parse some JavaScript, yielding a *ast.Program and/or an ErrorList
    program, err := parser.ParseFile(nil, filename, src, 0)

Warning

The parser and AST interfaces are still works-in-progress (particularly where
node types are concerned) and may change in the future.

*/
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
	str       string
	length    int
	base      int

	chr       rune // The current character
	chrOffset int  // The offset of current character
	offset    int  // The offset after current character (may be greater than 1)

	idx       file.Idx    // The index of token
	token     token.Token // The token
	literal   string      // The literal of the token, if any
	isKeyword bool

	scope *Scope

	insertSemicolon   bool // If we see a newline, then insert an implicit semicolon
	implicitSemicolon bool // An implicit semicolon exists

	arrowFunctionMode  bool // When trying to parse possible arrow fn parameters
	patternBindingMode bool // When trying to parse possible binding pattern

	jsxTextParseFrom int

	errors ErrorList

	recover struct {
		// Scratch when trying to seek to the next statement, etc.
		idx   file.Idx
		count int
	}

	file *file.File
}

func _newParser(filename, src string, base int) *Parser {
	return &Parser{
		chr:    ' ', // This is set so we can start scanning by skipping whitespace
		str:    src,
		length: len(src),
		base:   base,
		file:   file.NewFile(filename, src, base),
	}
}

func newParser(filename, src string) *Parser {
	return _newParser(filename, src, 1)
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

		parser := _newParser(filename, str, base)
		return parser.parse()
	}
}

// ParseFunction parses a given parameter list and body as a function and returns the
// corresponding ast.FunctionLiteral node.
//
// The parameter list, if any, should be a comma-separated list of identifiers.
//
func ParseFunction(parameterList, body string) (*ast.FunctionLiteral, error) {

	src := "(function(" + parameterList + ") {\n" + body + "\n})"

	parser := _newParser("", src, 1)
	program, err := parser.parse()
	if err != nil {
		return nil, err
	}

	return program.Body[0].(*ast.ExpressionStatement).Expression.(*ast.FunctionLiteral), nil
}

func (p *Parser) slice(idx0, idx1 file.Idx) string {
	from := int(idx0) - p.base
	to := int(idx1) - p.base
	if from >= 0 && to <= len(p.str) {
		return p.str[from:to]
	}

	return ""
}

func (p *Parser) parse() (*ast.Program, error) {
	p.next()
	program := p.parseProgram()
	if false {
		p.errors.Sort()
	}
	return program, p.errors.Err()
}

func (p *Parser) now() (token.Token, string, file.Idx) {
	return p.token, p.literal, p.idx
}

func (p *Parser) next() (idx file.Idx) {
	idx = p.idx
	p.token, p.literal, p.idx = p.scan()
	return
}

func (p *Parser) optionalSemicolon() {
	if p.is(token.SEMICOLON) {
		p.next()
		return
	}

	if p.implicitSemicolon {
		p.implicitSemicolon = false
		return
	}

	if !p.is(token.EOF) && !p.is(token.RIGHT_BRACE) {
		p.consumeExpected(token.SEMICOLON)
	}
}

func (p *Parser) semicolon() {
	if !p.is(token.RIGHT_PARENTHESIS) && !p.is(token.RIGHT_BRACE) {
		if p.implicitSemicolon {
			p.implicitSemicolon = false
			return
		}

		p.consumeExpected(token.SEMICOLON)
	}
}

func (p *Parser) idxOf(offset int) file.Idx {
	return file.Idx(p.base + offset)
}

func (p *Parser) is(value token.Token) bool {
	return p.token == value
}

func (p *Parser) isIdentifierOrKeyword() bool {
	return p.is(token.IDENTIFIER) || p.isKeyword
}

func (p *Parser) until(value token.Token) bool {
	return p.token != value && p.token != token.EOF
}

func (p *Parser) unexpectedToken() {
	p.errorUnexpectedToken(p.token)
}

func (p *Parser) consumeExpected(value token.Token) file.Idx {
	idx := p.idx
	if p.token != value {
		p.unexpectedToken()
	}
	p.next()
	return idx
}

func (p *Parser) consumePossible(value token.Token) file.Idx {
	idx := p.idx

	if p.token == value {
		p.next()
	}

	return idx
}

func (p *Parser) shouldBe(value token.Token) {
	if p.token != value {
		p.errorUnexpectedToken(p.token)
	}
}

func lineCount(str string) (int, int) {
	line, last := 0, -1
	pair := false
	for index, chr := range str {
		switch chr {
		case '\r':
			line += 1
			last = index
			pair = true
			continue
		case '\n':
			if !pair {
				line += 1
			}
			last = index
		case '\u2028', '\u2029':
			line += 1
			last = index + 2
		}
		pair = false
	}
	return line, last
}

func (p *Parser) position(idx file.Idx) file.Position {
	position := file.Position{}
	offset := int(idx) - p.base
	str := p.str[:offset]
	position.Filename = p.file.Name()
	line, last := lineCount(str)
	position.Line = 1 + line
	if last >= 0 {
		position.Column = offset - last
	} else {
		position.Column = 1 + len(str)
	}

	return position
}
