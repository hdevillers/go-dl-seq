package main

import (
	"flag"
	"fmt"

	"github.com/hdevillers/go-dl-seq/seqio"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	// Retrieve argument values
	input := flag.String("input", "STDIN", "Input fasta file")
	format := flag.String("format", "fasta", "Input format.")
	gunzip := flag.Bool("gunzip", false, "Input must is compressed.")
	flag.Parse()

	if *input == "" {
		panic("You must provide an input fasta file.")
	}

	// Open sequence file
	seqIn := seqio.NewReader(*input, *format, *gunzip)
	seqIn.CheckPanic()
	defer seqIn.Close()

	for seqIn.Next() {
		seqIn.CheckPanic()
		s := seqIn.Seq()

		fmt.Printf("%s\t%d\n", s.Id, s.Length())
	}
}
