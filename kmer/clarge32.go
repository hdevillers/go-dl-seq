package kmer

import (
	"errors"
	"fmt"
	"math"
	"sort"
)

type clarge32 struct {
	Km *Kmer32      // Kmer manager
	C  []uint64     // Kmer counter
	I  *Kmer32Slice // Kmer id ordered
	R  *Kmer32Slice // Raw Kmer list
	F  bool         // Finished indicator

}

func NewClarge32(k int) *clarge32 {
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
	j := 1
	for i := int(cl.Km.K); i < len(b); i++ {
		cl.Km.AddBase(b[i])
		cl.R.W[ci+j] = cl.Km.W
		j++
	}
}

func (cl *clarge32) Finish() {
	if cl.F {
		panic(errors.New("[KMER LARGE COUNTER]: Counter already finished."))
	}

	// First: sort the raw list of Kmers
	sort.Sort(cl.R)

	// Count ordered Kmers
	nK := cl.R.Len()
	tmpC := make([]int, nK)
	tmpI := make([]uint32, nK)
	i := 0
	s := 0
	for i < nK {
		// Count the number of consecutive identical Kmers
		tmpC[s] = 1
		tmpI[s] = cl.R.W[i]
		j := i + tmpC[s]
		for j < nK && cl.R.Equal(i, j) {
			tmpC[s]++
			j++
		}

		// Continue
		i += tmpC[s]
		s++
	}

	// Store results in the counter
	cl.I.Extend(s)
	cl.C = make([]uint64, s)
	for i = 0; i < s; i++ {
		cl.I.W[i] = tmpI[i]
		cl.C[i] = uint64(tmpC[i])
	}

	cl.F = true
}

func (cl *clarge32) PrintAll() {
	if !cl.F {
		panic(errors.New("[KMER SMALL COUNTER]: Before printing counted values, you must call the Finish method."))
	}
	k := int(cl.Km.K)
	if k > MaxKPrintAll {
		panic(errors.New("[KMER LARGE COUNTER]: K value to large to print all possibilities."))
	}
	n := math.Pow(4.0, float64(k))
	j := 0
	J := cl.I.Len()
	for i := 0; i < int(n); i++ {
		if j < J && cl.I.W[j] == uint32(i) {
			fmt.Printf("%s\t%d\n", Kmer32String(cl.I.W[j], k), cl.C[j])
			j++
		} else {
			fmt.Printf("%s\t%d\n", Kmer32String(uint32(i), k), 0)
		}
	}
}

func (cl *clarge32) Print() {
	if !cl.F {
		panic(errors.New("[KMER SMALL COUNTER]: Before printing counted values, you must call the Finish method."))
	}
	k := int(cl.Km.K)
	for i := 0; i < cl.I.Len(); i++ {
		fmt.Printf("%s\t%d\n", Kmer32String(cl.I.W[i], k), cl.C[i])
	}
}
