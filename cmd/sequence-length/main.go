package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/hdevillers/go-dl-seq/seqio/fasta"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	// Retrieve argument values
	input := flag.String("input", "", "Input fasta file")
	flag.Parse()

	if *input == "" {
		panic("You must provide an input fasta file.")
	}

	// Open fasta file
	f, err := os.Open(*input)
	check(err)
	defer f.Close()

	// Fasta Reader
	fio := fasta.NewReader(f)

	for !fio.IsEOF {
		s, err := fio.Read()
		check(err)

		fmt.Printf("%s\t%d\n", s.Id, s.Length())
	}
}
