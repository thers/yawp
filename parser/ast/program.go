package ast

import (
	"yawp/ids"
	"yawp/parser/file"
)

type ModuleAdditions struct {
	ObjectOmit bool
}

type Module struct {
	File *file.File
	Body []Statement

	Ids       *ids.Ids
	Additions ModuleAdditions
}

func (m *Module) GetLoc() *file.Loc {
	return m.Body[0].GetLoc().Add(m.Body[len(m.Body)-1].GetLoc())
}

func (m *Module) Visit(visitor Visitor) {
	m.Body = visitor.Body(m.Body)
}
