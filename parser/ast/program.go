package ast

import (
	"github.com/go-sourcemap/sourcemap"
	"yawp/parser/file"
)

type Program struct {
	Body []Statement

	DeclarationList []Declaration

	File *file.File

	SourceMap *sourcemap.Consumer
}

func (p *Program) GetLoc() *file.Loc {
	return p.Body[0].GetLoc().Add(p.Body[len(p.Body)-1].GetLoc())
}
