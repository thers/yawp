package ast

import (
	"yawp/parser/file"
)

// All nodes implement the INode interface.
type INode interface {
	GetLoc() *file.Loc
	GetNode() *Node
}

type Node struct {
	Loc  *file.Loc
	Flag Flags
}

func (n *Node) GetLoc() *file.Loc { return n.Loc }

func (n *Node) Copy() Node {
	return Node{
		Loc:  n.Loc.Copy(),
		Flag: n.Flag,
	}
}
