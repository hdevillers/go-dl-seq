package kmer

import (
	"errors"
	"fmt"
	"math"
)

type csmall struct {
	Km *Kmer32    // Kmer manager
	C  [][]uint64 // Kmer counter
	F  bool       // Finished indicator
	L  int        // Counter indice (multiple channel)
}

func NewCsmall(k int, l int) *csmall {
	// Do not instantiate for large k value
	if k > MaxKSmall {
		panic(errors.New("[KMER SMALL COUNTER]: Too large K value."))
	}

	// Initialize
	var cs csmall
	cs.Km = NewKmer32(k)
	n := math.Pow(4.0, float64(k))
	cs.C = make([][]uint64, l)
	for i := 0; i < l; i++ {
		cs.C[i] = make([]uint64, int64(n))
	}
	cs.F = false
	cs.L = 0

	// Return the counter
	return &cs
}

func (cs *csmall) Count(b []byte) {
	// Initialize the kmer manager
	cs.Km.Init(b)
	cs.F = false

	// Count the first word
	cs.C[cs.L][cs.Km.W]++

	// Count the following words
	for i := int(cs.Km.K); i < len(b); i++ {
		cs.Km.AddBase(b[i])
		cs.C[cs.L][cs.Km.W]++
	}
}

func (cs *csmall) NextChannel() {
	cs.L++
}

func (cs *csmall) Finish() {
	// In small counter, just add a channel
	cs.NextChannel()
	cs.F = true
}

func (cs *csmall) PrintAll() {
	if !cs.F {
		panic(errors.New("[KMER SMALL COUNTER]: Before printing counted values, you must call the Finish method."))
	}
	k := int(cs.Km.K)
	for i := 0; i < len(cs.C[0]); i++ {
		fmt.Printf("%s", Kmer32String(uint32(i), k))
		for j := 0; j < cs.L; j++ {
			fmt.Printf("\t%d", cs.C[j][i])
		}
		fmt.Printf("\n")
	}
}

func (cs *csmall) Print() {
	if !cs.F {
		panic(errors.New("[KMER SMALL COUNTER]: Before printing counted values, you must call the Finish method."))
	}
	k := int(cs.Km.K)
	for i := 0; i < len(cs.C[0]); i++ {
		sum := uint64(0)
		for j := 0; j < cs.L; j++ {
			sum += cs.C[j][i]
		}
		if sum > 0 {
			fmt.Printf("%s", Kmer32String(uint32(i), k))
			for j := 0; j < cs.L; j++ {
				fmt.Printf("\t%d", cs.C[j][i])
			}
			fmt.Printf("\n")
		}
	}
}
