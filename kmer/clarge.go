package kmer

import (
	"errors"
)

type clarge struct {
	Km KmerManager // Kmer manager
	C  []uint64    // Kmer counter
	I  KmerSlice   // Kmer id ordered
	R  KmerSlice   // Raw Kmer list
	F  bool        // Finished indicator

}

func NewClarge(k int) *clarge {
	// Do not instantiate if K is to large
	if k > MaxKAbsolute {
		panic(errors.New("[KMER LARGE COUNTER]: Too large K value."))
	}

	// Initialize
	var cl clarge
	if k <= MaxK32Bits {
		// Can use 32 uint to count
		cl.Km = NewKmer32(k)
		cl.C = make([]uint64, 0)
		tmp := make([]uint32, 0)
		cl.I = tmp
	}
	cl.F = false

	return &cl
}
