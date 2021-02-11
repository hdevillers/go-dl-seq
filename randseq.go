package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/hdevillers/go-dl-seq/seq"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	// Retrieve argument values
	seqLen := flag.Int("len", 200, "Length of required sequences.")
	input := flag.String("input", "", "Input fasta file")
	flag.Parse()

	if *input == "" {
		panic("You must provide an input fasta file.")
	}

	// Open fasta file
	f, err := os.Open(*input)
	check(err)
	defer f.Close()

	// New seq
	id := "MonId"
	sequence := []byte{'A', 'C', 'G', 'A'}
	s := seq.NewSeq(id, sequence)
	fmt.Printf("Sequence length is %d.", s.Length())

	r := bufio.NewReader(f)
	line, err := r.ReadBytes('\n')
	if err == io.EOF {
		panic("Empty file!")
	}
	if line[0] == '>' {
		fmt.Println("Good fasta header.")
		fmt.Println(string(line[1:]))
	}

	fmt.Println("Sequence length: ", *seqLen, " nuc.")
}
