package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/hdevillers/go-seq/seq"

	"github.com/hdevillers/go-seq/seqio"
	"github.com/hdevillers/go-seq/utils"
)

func main() {
	input := flag.String("input", "", "Input sequence file to split.")
	output := flag.String("output", "dataset_split", "Output base name.")
	format := flag.String("format", "fasta", "Format of input and output.")
	gz := flag.Bool("gz", false, "Input and outputs are compressed (gz).")
	pTrain := flag.Float64("ptrain", 0.7, "Proportion of sequences in the training set.")
	pTest := flag.Float64("ptest", 0.2, "Proportion of sequences in the test set.")
	pValid := flag.Float64("pvalid", 0.1, "Propotion of sequences in the validation set.")
	userSeed := flag.Int64("seed", 0, "Provide the random seed.")
	flag.Parse()

	// Argument check
	if *input == "" {
		panic("You must provide an input sequence file to split.")
	}

	// Sum of proportions must be 1.0
	pSum := *pTrain + *pTest + *pValid
	if math.Abs(pSum-1.0) >= 1e-6 {
		panic("Sum of the proportion of the training, validation and test sets must be equal to 1.")
	}

	// Load input sequences
	var seqs []seq.Seq
	nseqs := utils.LoadSeqInArray(*input, *format, &seqs)
	if nseqs == 0 {
		panic("No sequence found in the provided input.")
	}

	// Shuffle the sequences
	seed := *userSeed
	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	fmt.Printf("Used random seed: %d\n", seed)
	seeder := rand.NewSource(seed)
	random := rand.New(seeder)
	random.Shuffle(nseqs, func(i, j int) {
		seqs[i], seqs[j] = seqs[j], seqs[i]
	})

	// Compute the number of sequences for each datasets
	nTrain := int(*pTrain * float64(nseqs))
	nValid := int(*pValid * float64(nseqs))
	nTest := nseqs - nTrain - nValid

	// Prepare output files
	seqTrain := seqio.NewWriter(fmt.Sprintf("%s.train.%s", *output, *format), *format, *gz)
	defer seqTrain.Close()
	seqValid := seqio.NewWriter(fmt.Sprintf("%s.valid.%s", *output, *format), *format, *gz)
	defer seqValid.Close()
	seqTest := seqio.NewWriter(fmt.Sprintf("%s.test.%s", *output, *format), *format, *gz)
	defer seqTest.Close()

	// Write out
	for i := 0; i < nTrain; i++ {
		seqTrain.Write(seqs[i])
	}
	for i := 0; i < nValid; i++ {
		seqValid.Write(seqs[nTrain+i])
	}
	for i := 0; i < nTest; i++ {
		seqValid.Write(seqs[nTrain+nValid+i])
	}

	// Print a summary
	fmt.Printf("Wrote %d sequences in the training set.\n", nTrain)
	fmt.Printf("Wrote %d sequences in the validation set.\n", nValid)
	fmt.Printf("Wrote %d sequences in the test set.\n", nTest)
}
