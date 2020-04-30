package ast

import "yawp/parser/file"

type Comment struct {
	From file.Idx
	To   file.Idx

	String string
}
