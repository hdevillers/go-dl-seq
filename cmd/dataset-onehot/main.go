package main

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/hdevillers/go-dl-seq/dataset"
	"github.com/hdevillers/go-dl-seq/encoding"
	"github.com/hdevillers/go-dl-seq/seq"
	"github.com/hdevillers/go-dl-seq/seqio"
	"gonum.org/v1/hdf5"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func fastaToDataset(input string, format string, h5 *hdf5.File, dtype string) (int, int) {
	// Initialize variables
	var seqs []seq.Seq
	var labels []string
	var y []uint16

	// Open input fasta file an scan it
	slen := 0
	seqIn := seqio.NewReader(input, format)
	seqIn.CheckPanic()
	defer seqIn.Close()
	for seqIn.Next() {
		seqIn.CheckPanic()
		// Initialize the encoder
		s := seqIn.Seq()
		seqs = append(seqs, s)
		labels = append(labels, s.Id)
		tmpY, err := strconv.Atoi(s.Desc)
		check(err)
		y = append(y, uint16(tmpY))

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
	x := make([]uint16, nseq*slen*4)
	for i, s := range seqs {
		oho := encoding.NewOneHot(s.Sequence)
		oho.Compute()
		oho.CheckPanic()
		at := slen * 4 * i
		copy(x[at:], oho.GetOho())
	}

	// Add the data set to the h5 file
	err := dataset.AddOhoDataset(h5, dtype, x, y, labels, nseq, slen)
	check(err)

	return slen, nseq
}

func main() {
	// Retrieve argument values
	training := flag.String("train", "", "Input training sequence file.")
	testing := flag.String("test", "", "Input testing sequence file.")
	validation := flag.String("valid", "", "Input validation sequence file.")
	format := flag.String("format", "fasta", "Input format.")
	output := flag.String("output", "dataset.h5", "Output dataset file (h5).")
	flag.Parse()

	if *training == "" {
		panic("You must provide an input training file.")
	}

	// Create the H5 file
	h5, err := hdf5.CreateFile(*output, hdf5.F_ACC_TRUNC)
	check(err)

	// Generate the training dataset
	ltrain, ntrain := fastaToDataset(*training, *format, h5, "training")
	fmt.Printf("Training set: %d sequences of %d bases.\n", ntrain, ltrain)

	// Generate the testing dataset
	if *testing != "" {
		ltest, ntest := fastaToDataset(*testing, *format, h5, "testing")
		if ltest != ltrain {
			panic("Training and testing sequences must have the same length.")
		}
		fmt.Printf("Testing set: %d sequences of %d bases.\n", ntest, ltest)
	} else {
		fmt.Println("No testing set.")
	}

	// Generate the validation dataset
	if *validation != "" {
		lvalid, nvalid := fastaToDataset(*validation, *format, h5, "validation")
		if lvalid != ltrain {
			panic("Training and validation sequences must have the same length.")
		}
		fmt.Printf("Validation set: %d sequences of %d bases.\n", nvalid, lvalid)
	} else {
		fmt.Println("No validation set.")
	}

	h5.Close()
}
