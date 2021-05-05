package main

import (
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

	kmerCounter := kmer.NewCsmall(*k)

	for seqIn.Next() {
		seqIn.CheckPanic()
		s := seqIn.Seq()

		kmerCounter.Count(s.Sequence)
	}

	kmerCounter.PrintAll()
}
