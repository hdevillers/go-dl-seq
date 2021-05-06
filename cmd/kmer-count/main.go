package main

import (
	"errors"
	"flag"

	"github.com/hdevillers/go-dl-seq/kmer"
	"github.com/hdevillers/go-dl-seq/seqio"
)

func main() {
	k := flag.Int("k", 4, "K value.")
	i := flag.String("i", "", "Input sequence file.")
	f := flag.String("f", "fasta", "Input sequence format.")
	d := flag.Bool("d", false, "Decompress the input (gz).")
	flag.Parse()

	if *i == "" {
		panic("You must provide an input fasta file.")
	}

	seqIn := seqio.NewReader(*i, *f, *d)

	// Determine the type of counter
	var kmerCounter kmer.KmerCounter
	if *k <= kmer.MaxKSmall {
		kmerCounter = kmer.NewCsmall(*k)
	} else if *k <= kmer.MaxK32Bits {
		kmerCounter = kmer.NewClarge32(*k)
	} else {
		panic(errors.New("K value is too large."))
	}

	// Count Kmer in all input sequences
	for seqIn.Next() {
		seqIn.CheckPanic()
		s := seqIn.Seq()
		kmerCounter.Count(s.Sequence)
	}

	// Terminate counter
	kmerCounter.Finish()

	// Print out counted value
	kmerCounter.Print()
}
