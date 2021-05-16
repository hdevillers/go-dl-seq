package main

import (
	"flag"
	"fmt"
	"sort"

	"github.com/hdevillers/go-dl-seq/seqio"
)

type counter struct {
	id int
	n  int
	K  int
	B  []uint32 // Converter
	F  int
	I  []uint32 // Kmer IDs
	C  []uint32 // Kmer counts
	R  int
}

type counters struct {
	counter []*counter
}

func newCounters(k int) *counters {
	var c counters
	c.counter = make([]*counter, k)
	for i := 0; i < k; i++ {
		c.counter[i] = newCounter(i, k)
	}
	return &c
}

func newCounter(id int, k int) *counter {
	var c counter
	c.id = id
	c.K = k
	c.B = make([]uint32, 256)
	c.F = (16 - k + 1) * 2
	c.R = (16 - k) * 2

	// Set base values
	c.B['C'] = uint32(1)
	c.B['c'] = uint32(1)
	c.B['G'] = uint32(2)
	c.B['g'] = uint32(2)
	c.B['T'] = uint32(3)
	c.B['t'] = uint32(3)

	return &c
}

func (c *counter) count(seqChan chan []byte, couChan chan int) {
	rawList := make([]uint32, 0)
	nw := 0 // Number of words (in rawList)
	for seq := range seqChan {
		l := len(seq)

		// Extend rawList
		rawList = append(rawList, make([]uint32, l-c.K+1)...)

		// Init the first words
		w := uint32(0) // Runing word
		for i := 0; i < c.K; i++ {
			w = (w << 2) | c.B[seq[i]]
		}

		// Add the first word
		rawList[0] = w
		nw++

		// Continue to enumerate words
		for i := c.K; i < l; i++ {
			w = (w<<c.F)>>c.R | c.B[seq[i]]
			rawList[nw] = w
			nw++
		}
	}

	c.n = nw

	/*
		Sort
	*/
	sort.Slice(rawList, func(i, j int) bool {
		return rawList[i] < rawList[i]
	})

	/*
		Count
	*/
	// Initialize the counter container
	c.I = make([]uint32, nw)
	c.C = make([]uint32, nw)

	// Read words
	i := 0
	stored := 0
	for i < nw {
		c.I[stored] = rawList[i]
		cou := 1
		j := i + 1
		for j < nw && rawList[i] == rawList[j] {
			cou++
			j++
		}
		c.C[stored] = uint32(cou)
		i += cou
		stored++
	}

	// Number of counted words
	c.n = stored

	couChan <- c.id
}

func (c *counters) findPairs(paiChan chan int, max int) {
	n := 0
	i := 0
	for n < max && i < len(c.counter) {
		if c.counter[i] != nil {
			paiChan <- i
			n++
		}
		i++
	}
	if n < max {
		// Should not occure!
		panic("Merging thread issue!")
	}
}

func (c *counters) merge(paiChan chan int, merChan chan int) {
	i := <-paiChan
	j := <-paiChan
	// Merging counters i and j...
	fmt.Println("Merging counter ", i, " with counter ", j)

	// Erasing counter j
	c.counter[j] = nil

	merChan <- i
}

func main() {
	input := flag.String("i", "", "Input sequence.")
	k := flag.Int("k", 4, "Kmer values.")
	t := flag.Int("t", 4, "Number of threads.")
	flag.Parse()

	if *input == "" {
		panic("No input!")
	}

	// Number of threads
	threads := *t

	// Sequence channel
	seqChan := make(chan []byte, threads)
	couChan := make(chan int)
	paiChan := make(chan int)
	merChan := make(chan int)

	// Init counters
	counters := newCounters(*k)
	for i := 0; i < threads; i++ {
		go counters.counter[i].count(seqChan, couChan)
	}

	// SeqIO input sequences
	seqIn := seqio.NewReader(*input, "fasta", false)
	defer seqIn.Close()

	for seqIn.Next() {
		seqIn.CheckPanic()
		s := seqIn.Seq()
		seqChan <- s.Sequence
	}
	close(seqChan)

	// Wait until all counters are done
	for i := 0; i < threads; i++ {
		<-couChan
	}

	// Merging counters
	nc := threads // Number of counters
	nm := nc / 2  // Number of merging process
	rm := nc % 2  // Number of unmerged counter
	for nc > 1 {
		// merging go routine
		for i := 0; i < nm; i++ {
			go counters.merge(paiChan, merChan)
		}

		// through pairs
		counters.findPairs(paiChan, 2*nm)

		// Wait the merged counters
		for i := 0; i < nm; i++ {
			<-merChan
		}

		// refine numbers
		nc = nm + rm
		nm = nc / 2
		rm = nc % 2
	}

	for i := 0; i < threads; i++ {
		if counters.counter[i] != nil {
			fmt.Println("Count ", counters.counter[i].id, " number of words: ", counters.counter[i].n)
		}
	}

}
