// Package file encapsulates the file abstractions used by the ast & parser.
//
package file

import (
	"fmt"
	"strings"
)

// Start is a compact encoding of a source position within a file set.
// It can be converted into a Pos for a more convenient, but much
// larger, representation.
type Idx int

type Loc struct {
	From Idx
	To   Idx

	Line int
	Col  int
}

func (l *Loc) End(offset Idx) *Loc {
	l.To = offset

	return l
}

func (l *Loc) Add(loc *Loc) *Loc {
	c := l.Copy()
	c.To = loc.To

	return c
}

func (l *Loc) Copy() *Loc {
	if l == nil {
		return &Loc{}
	}

	return &Loc{
		From: l.From,
		To:   l.To,
		Line: l.Line,
		Col:  l.Col,
	}
}

// Pos describes an arbitrary source position
// including the filename, line, and column location.
type Pos struct {
	Loc    int // The src offset
	Line   int // The line number, starting at 1
	Column int // The column number, starting at 1 (The character count)
}

// A Pos is valid if the line number is > 0.

func (self *Pos) isValid() bool {
	return self.Line > 0
}

// String returns a string in one of several forms:
//
//	file:line:column    A valid position with filename
//	line:column         A valid position without filename
//	file                An invalid position with filename
//	-                   An invalid position without filename
//
func (self *Pos) String() string {
	str := ""
	if self.isValid() {
		if str != "" {
			str += ":"
		}
		str += fmt.Sprintf("%d:%d", self.Line, self.Column)
	}
	if str == "" {
		str = "-"
	}
	return str
}

// FileSet

// A FileSet represents a set of source files.
type FileSet struct {
	files []*File
	last  *File
}

// AddFile adds a new file with the given filename and src.
//
// This an internal method, but exported for cross-package use.
func (self *FileSet) AddFile(filename, src string) int {
	base := self.nextBase()
	file := &File{
		name: filename,
		src:  src,
		base: base,
	}
	self.files = append(self.files, file)
	self.last = file
	return base
}

func (self *FileSet) nextBase() int {
	if self.last == nil {
		return 1
	}
	return self.last.base + len(self.last.src) + 1
}

func (self *FileSet) File(idx Idx) *File {
	for _, file := range self.files {
		if idx <= Idx(file.base+len(file.src)) {
			return file
		}
	}
	return nil
}

// Pos converts an Start in the FileSet into a Pos.
func (self *FileSet) Position(idx Idx) *Pos {
	position := &Pos{}
	for _, file := range self.files {
		if idx <= Idx(file.base+len(file.src)) {
			offset := int(idx) - file.base
			src := file.src[:offset]
			position.Loc = offset
			position.Line = 1 + strings.Count(src, "\n")
			if index := strings.LastIndex(src, "\n"); index >= 0 {
				position.Column = offset - index
			} else {
				position.Column = 1 + len(src)
			}
		}
	}
	return position
}

type File struct {
	name string
	src  string
	base int // This will always be 1 or greater
}

func NewFile(filename, src string) *File {
	return &File{
		name: filename,
		src:  src,
	}
}

func (fl *File) Name() string {
	return fl.name
}

func (fl *File) Source() string {
	return fl.src
}

func (fl *File) Base() int {
	return fl.base
}
