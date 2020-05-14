package ast

type Flag uint64

func (f Flag) Add(index uint64) Flag {
	return f | (1 << (index - 1))
}

func (f Flag) Delete(index uint64) Flag {
	return f ^ (1 << (index - 1))
}

func (f Flag) Has(flags Flag) bool {
	return (f & flags) > 0
}
