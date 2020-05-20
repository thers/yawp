package ast

type Flags uint64

func (f Flags) Add(flag Flags) Flags {
	return f | flag
}

func (f Flags) Delete(flag Flags) Flags {
	return f ^ flag
}

func (f Flags) Has(flags Flags) bool {
	return (f & flags) > 0
}
