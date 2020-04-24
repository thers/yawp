package ast

import (
	"yawp/ids"
	"yawp/parser/file"
)

type Module struct {
	File *file.File
	Body []Statement
	Ids  *ids.Ids
}

func (m *Module) GetLoc() *file.Loc {
	return m.Body[0].GetLoc().Add(m.Body[len(m.Body)-1].GetLoc())
}

func (m *Module) Visit(visitor Visitor) {
	for index, stmt := range m.Body {
		m.Body[index] = visitor.Statement(stmt)
	}
}
