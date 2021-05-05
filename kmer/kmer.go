package kmer

const (
	MaxKSmall    int = 12
	MaxKAbsolute int = 15
)

type KmerManager interface {
	Init([]byte)
	AddBase(byte)
}

type KmerSlice interface {
	Len() int
	Less(int, int)
	Swap(int, int)
}

type KmerCounter interface {
	Count([]byte)
	Finish()
}
