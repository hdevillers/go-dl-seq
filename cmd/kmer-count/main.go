package main

import (
	"errors"
	"flag"
	"log"
	"os"

	"github.com/hdevillers/go-dl-seq/kmer"
	"github.com/hdevillers/go-dl-seq/seqio"
)

/* Define an input type to allow multiple input files */
type inputFlags []string

func (i *inputFlags) String() string {
	return "hello, world\n"
}

func (i *inputFlags) Set(s string) error {
	*i = append(*i, s)
	return nil
}

var input inputFlags

func main() {
	k := flag.Int("k", 4, "K value.")
	flag.Var(&input, "i", "Input sequence file(s).")
	f := flag.String("f", "fasta", "Input sequence format.")
	d := flag.Bool("d", false, "Decompress the input (gz).")
	a := flag.Bool("a", false, "Print all Kmers, including zero-count.")
	g := flag.Bool("g", false, "Group multiple file in a single counter.")
	flag.Parse()

	logger := log.New(os.Stderr, "DEBUG: ", log.Lmicroseconds)

	if len(input) == 0 {
		panic("You must provide at one input fasta file.")
	}

	if *a {
		if *k > kmer.MaxKPrintAll {
			panic("K value is too large to print all possible Kmers.")
		}
	}

	// Number of requiered channel
	nc := len(input)
	if *g {
		nc = 1
	}

	// Determine the type of counter
	var kmerCounter kmer.KmerCounter
	logger.Print("Initializing counter...")
	if *k <= kmer.MaxKSmall {
		kmerCounter = kmer.NewCsmall(*k, nc)
	} else if *k <= kmer.MaxK32Bits {
		kmerCounter = kmer.NewClarge32(*k, nc)
	} else {
		panic(errors.New("K value is too large."))
	}

	logger.Print("Start reading sequences...")
	for i := 0; i < len(input); i++ {
		seqIn := seqio.NewReader(input[i], *f, *d)

		// Count Kmer in all input sequences
		for seqIn.Next() {
			seqIn.CheckPanic()
			s := seqIn.Seq()
			kmerCounter.Count(s.Sequence)
		}

		// Finish and add a channel if not grouped
		if !*g {
			kmerCounter.Finish()
		}
	}

	// Finish the overall count if grouped
	logger.Print("Start sorting and counting...")
	if *g {
		kmerCounter.Finish()
	}

	logger.Print("Start writing out...")
	// Print out counted value
	if *a {
		kmerCounter.PrintAll()
	} else {
		kmerCounter.Print()
	}
	logger.Print("Finished.")
}
