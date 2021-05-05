package kmer

const (
	MaxKSmall    int = 12
	MaxK32Bits   int = 15
	MaxK64Bits   int = 31
	MaxKAbsolute int = 31
)

type KmerCounter interface {
	Count([]byte)
	Finish()
}
