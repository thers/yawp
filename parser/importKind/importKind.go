package importKind

type ImportKind int

const (
	VALUE ImportKind = iota
	TYPE
	TYPEOF
)
