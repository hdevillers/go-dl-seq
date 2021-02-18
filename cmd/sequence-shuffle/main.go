package main

import (
	"flag"
	"math/rand"
	"time"

	"github.com/hdevillers/go-dl-seq/seqio"
	"github.com/hdevillers/go-dl-seq/seq"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	// Retrieve argument values
	input := flag.String("input", "", "Input sequence file.")
	format := flag.String("format", "fasta", "Input/output format.")
	output := flag.String("output", "shuffled.fasta", "Output sequence file.")
	flag.Parse()

	if *input == "" {
		panic("You must provide an input fasta file.")
	}

	// Setup random seed
	seeder := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seeder)

	// Read input sequences
	var seqs []seq.Seq
	seqIn := seqio.NewReader(*input, *format)
	seqIn.CheckPanic()
	defer seqIn.Close()
	for seqIn.Next() {
		seqIn.CheckPanic()
		seqs = append(seqs, seqIn.Seq())
	}

	// Shuffle the slice of sequences
	random.Shuffle(len(seqs), func(i,j int) {
		seqs[i], seqs[j] = seqs[j], seqs[i]
	})

	// Save shuffled sequences in output
	seqOut := seqio.NewWriter(*output, *format)
	seqOut.CheckPanic()
	defer seqOut.Close()

	for _,s := range seqs {
		seqOut.Write(s)
	}
}
