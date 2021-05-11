package kmer

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
)

type clarge32 struct {
	Km *Kmer32      // Kmer manager
	C  [][]uint64   // Kmer counter
	I  *Kmer32Slice // Kmer id ordered
	R  *Kmer32Slice // Raw Kmer list
	Ch int          // Channel (multiple counter)
	F  bool         // Finished indicator

}

func NewClarge32(k int, ch int) *clarge32 {
	// Do not instantiate if K is to large
	if k > MaxK32Bits {
		panic(errors.New("[KMER LARGE COUNTER]: Too large K value."))
	}

	// Initialize
	var cl clarge32
	cl.Km = NewKmer32(k)
	cl.C = make([][]uint64, ch)
	for i := 0; i < ch; i++ {
		cl.C[i] = make([]uint64, 0)
	}
	cl.I = NewKmer32Slice()
	cl.R = NewKmer32Slice()
	cl.Ch = 0
	cl.F = false

	return &cl
}

func (cl *clarge32) Count(b []byte) {
	// Get the current counter length
	ci := cl.R.Len()
	cl.F = false

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

	// Initialize temp variables
	nK := cl.R.Len()
	nI := cl.I.Len()
	iI := 0
	tmpC := make([][]uint64, cl.Ch+1)
	for i := 0; i <= cl.Ch; i++ {
		tmpC[i] = make([]uint64, nK+nI)
	}
	tmpI := make([]uint32, nK+nI)

	// Count Kmer and combine with previous counts
	i := 0
	s := 0
	for i < nK {
		// First, count the number of occurence of the current Kmer
		curI := cl.R.W[i]
		curC := 1
		j := i + 1
		for j < nK && cl.R.Equal(i, j) {
			curC++
			j++
		}

		// Store all the Kmer already counted I[] < curI
		for iI < nI && cl.I.W[iI] < curI {
			// There are previously counter Kmer, not found here
			tmpI[s] = cl.I.W[iI]
			for k := 0; k < cl.Ch; k++ {
				tmpC[k][s] = cl.C[k][iI]
			}
			tmpC[cl.Ch][s] = 0 // NOTE: Useless...
			iI++               // Next previous entry
			s++                // Next new entry
		}
		// Check if the current I[] is equal to curI
		if iI < nI && cl.I.W[iI] == curI {
			tmpI[s] = cl.I.W[iI]
			for k := 0; k < cl.Ch; k++ {
				tmpC[k][s] = cl.C[k][iI]
			}
			tmpC[cl.Ch][s] = uint64(curC)
			iI++ // Next previous entry
			s++  // Next new entry
		} else {
			// This is a new Kmer!
			tmpI[s] = curI
			for k := 0; k < cl.Ch; k++ {
				tmpC[k][s] = 0
			}
			tmpC[cl.Ch][s] = uint64(curC)
			s++
		}

		// Continue to next current Kmer
		i += curC
	}

	// Copy temp variable into the counter
	cl.I = NewKmer32Slice()
	cl.I.Extend(s)                 // Reset name list
	cl.NextChannel()               // Add new channel
	cl.C = make([][]uint64, cl.Ch) // Reset counter
	for k := 0; k < cl.Ch; k++ {
		cl.C[k] = make([]uint64, s)
	}
	for j := 0; j < s; j++ {
		cl.I.W[j] = tmpI[j]
		for k := 0; k < cl.Ch; k++ {
			cl.C[k][j] = tmpC[k][j]
		}
	}

	cl.F = true
}

func (cl *clarge32) NextChannel() {
	// Reset the raw counter
	cl.R = NewKmer32Slice()
	cl.Ch++
}

func (cl *clarge32) PrintAll() {
	if !cl.F {
		panic(errors.New("[KMER LARGE COUNTER]: Before printing counted values, you must call the Finish method."))
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
			fmt.Printf("%s", Kmer32String(cl.I.W[j], k))
			for k := 0; k < cl.Ch; k++ {
				fmt.Printf("\t%d", cl.C[k][j])
			}
			j++
		} else {
			fmt.Printf("%s", Kmer32String(uint32(i), k))
			for k := 0; k < cl.Ch; k++ {
				fmt.Printf("\t%d", 0)
			}
		}
		fmt.Printf("\n")
	}
}

func (cl *clarge32) Print(output string) {
	if !cl.F {
		panic(errors.New("[KMER SMALL COUNTER]: Before printing counted values, you must call the Finish method."))
	}

	f, e := os.Create(output)
	if e != nil {
		panic(e)
	}
	defer f.Close()
	b := bufio.NewWriter(f)

	k := int(cl.Km.K)
	for i := 0; i < cl.I.Len(); i++ {
		b.WriteString(Kmer32String(cl.I.W[i], k))
		for j := 0; j < cl.Ch; j++ {
			b.WriteByte('\t')
			b.WriteString(strconv.FormatUint(cl.C[j][i], 10))
		}
		b.WriteByte('\n')
	}
	e = b.Flush()
	if e != nil {
		panic(e)
	}
}
