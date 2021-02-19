package main

import (
	"flag"
	"strconv"

	"gonum.org/v1/hdf5"
	"github.com/hdevillers/go-dl-seq/dataset"
	"github.com/hdevillers/go-dl-seq/seqio"
	"github.com/hdevillers/go-dl-seq/seq"
	"github.com/hdevillers/go-dl-seq/encoding"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	// Retrieve argument values
	input := flag.String("input", "", "Input sequence file.")
	format := flag.String("format", "fasta", "Input format.")
	output := flag.String("output", "dataset.h5", "Output dataset file (h5).")
	flag.Parse()

	if *input == "" {
		panic("You must provide an input fasta file.")
	}
	
	// Read input sequences
	var seqs []seq.Seq
	var trainLabels []string
	var trainY []uint16

	slen := 0
	seqIn := seqio.NewReader(*input, *format)
	seqIn.CheckPanic()
	defer seqIn.Close()
	for seqIn.Next() {
		seqIn.CheckPanic()
		// Initialize the encoder
		s := seqIn.Seq()
		seqs = append(seqs, s)
		trainLabels = append(trainLabels, s.Id)
		tmpY, err := strconv.Atoi(s.Desc)
		check(err)
		trainY = append(trainY, uint16(tmpY))

		// Throw an error if lengths differ between sequences
		if slen == 0 {
			slen = s.Length()
		} else {
			if slen != s.Length() {
				panic("With OneHot encoding, all sequences must have the same length.")
			}
		}
	}
	nseq := len(seqs)

	// Prepare OneHot coding data
	trainX := make([]uint16, nseq * slen * 4)
	for i,s := range seqs {
		oho := encoding.NewOneHot(s.Sequence)
		oho.Compute()
		oho.CheckPanic()
		at := slen * 4 * i
		copy(trainX[at:], oho.GetOho())
	}

	// Create the H5 file
	h5, err := hdf5.CreateFile(*output, hdf5.F_ACC_TRUNC)
	check(err)

	// Save training dataset
	dataset.AddOhoDataset(h5, "training", trainX, trainY, trainLabels, nseq, slen)

	
	h5.Close()
}
