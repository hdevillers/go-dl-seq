package kmer

import (
	"errors"
)

type clarge32 struct {
	Km *Kmer32      // Kmer manager
	C  []uint64     // Kmer counter
	I  *Kmer32Slice // Kmer id ordered
	R  *Kmer32Slice // Raw Kmer list
	F  bool         // Finished indicator

}

func NewClarge(k int) *clarge32 {
	// Do not instantiate if K is to large
	if k > MaxK32Bits {
		panic(errors.New("[KMER LARGE COUNTER]: Too large K value."))
	}

	// Initialize
	var cl clarge32
	cl.Km = NewKmer32(k)
	cl.C = make([]uint64, 0)
	cl.I = NewKmer32Slice()
	cl.R = NewKmer32Slice()
	cl.F = false

	return &cl
}

func (cl *clarge32) Count(b []byte) {
	// Get the current counter length
	ci := cl.R.Len()

	// Extend the counter
	cl.R.Extend(len(b) - int(cl.Km.K) + 1)

	// Initial
	cl.Km.Init(b)

	// Set the first word
	cl.R.W[ci] = cl.Km.W

	// Continue
	for i := int(cl.Km.K); i < len(b); i++ {
		cl.Km.AddBase(b[i])
		cl.R.W[ci+1] = cl.Km.W
	}
}
