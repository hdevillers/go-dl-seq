package main

import (
	"crypto/sha256"
	"flag"
	"fmt"

	"github.com/hdevillers/go-seq/seqio"
)

func main() {
	input := flag.String("input", "", "Input sequence file.")
	format := flag.String("format", "fasta", "Input format.")
	gz := flag.Bool("gz", false, "Input sequence is compressed (gz).")
	flag.Parse()

	// Check argument values
	if *input == "" {
		panic("You must provide an input sequence file.")
	}

	// Create data structure
	check := make(map[string]string)
	ndup := 0

	// Open sequence files
	seqIn := seqio.NewReader(*input, *format, *gz)
	seqIn.CheckPanic()
	defer seqIn.Close()

	// Scan each sequences
	for seqIn.Next() {
		seqIn.CheckPanic()
		seq := seqIn.Seq()
		key := fmt.Sprintf("%x", sha256.Sum256(seq.Sequence))
		val, test := check[key]
		if test {
			ndup++
			fmt.Printf("%s is identical to %s.\n", seq.Id, val)
		} else {
			check[key] = seq.Id
		}
	}

	if ndup == 0 {
		fmt.Println("No duplicated sequence found.")
	} else {
		fmt.Printf("%d duplicated sequences found.\n", ndup)
	}
}
