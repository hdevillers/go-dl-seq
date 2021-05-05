package kmer

import (
	"errors"
	"fmt"
	"math"
)

// TODO: Slice of id to be similar to large counter

type csmall struct {
	Km *Kmer32  // Kmer manager
	C  []uint64 // Kmer counter
	F  bool     // Finished indicator
}

func NewCsmall(k int) *csmall {
	// Do not instantiate for large k value
	if k > MaxKSmall {
		panic(errors.New("[KMER SMALL COUNTER]: Too large K value."))
	}

	// Initialize
	var cs csmall
	cs.Km = NewKmer32(k)
	n := math.Pow(4.0, float64(k))
	cs.C = make([]uint64, int64(n))
	cs.F = false

	// Return the counter
	return &cs
}

func (cs *csmall) Count(b []byte) {
	// Initialize the kmer manager
	cs.Km.Init(b)

	// Count the first word
	cs.C[cs.Km.W]++

	// Count the following words
	for i := int(cs.Km.K); i < len(b); i++ {
		cs.Km.AddBase(b[i])
		cs.C[cs.Km.W]++
	}
}

func (cs *csmall) Finish() {
	// In small counter, no thing to do
	cs.F = true
}

func (cs *csmall) PrintAll() {
	if !cs.F {
		panic(errors.New("[KMER SMALL COUNTER]: Before printing counted values, you must call the Finish method."))
	}
	for i := 0; i < len(cs.C); i++ {
		fmt.Printf("%s\t%d\n", Kmer32String(uint32(i), int(cs.Km.K)), cs.C[i])
	}
}
