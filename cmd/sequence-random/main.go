package main

import (
	"flag"
	"fmt"
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
	output := flag.String("output", "", "Output file name/path.")
	format := flag.String("format", "fasta", "Input format.")
	length := flag.Int("length", 200, "Required sequence length.")
	count  := flag.Int("n", 1, "Number of required sequence(s).")
	base   := flag.String("base", "RandSeq_", "Sequence ID base name.")
	seed   := flag.Int64("seed", 0, "Random seed initializer.")
	flag.Parse()

	if *output == "" {
		panic("You must provide an input fasta file.")
	}

	if *length <= 0 {
		panic("Sequence length must be greater than 0.")
	}

	if *count <= 0 {
		panic("The number of required sequence must be greater than 0.")
	}
	if *seed == 0 {
		// Initialize the seed with current time
		*seed = time.Now().UnixNano()
	}
	seeder := rand.NewSource(*seed)
	random := rand.New(seeder)

	fmt.Println("Used random seed:", *seed)

	// Open ouput file
	seqOut := seqio.NewWriter(*output, *format)
	seqOut.CheckPanic()
	defer seqOut.Close()

	// Very simple solution
	alpha := []byte{'A', 'C', 'G', 'T'}

	// Generate the required sequences
	for i:=0 ; i<*count ; i++ {
		// Create the new ID
		id := *base + fmt.Sprintf("%06d", i)

		// Create the new seq object
		seq := seq.NewSeq(id)

		// Add a sequence
		str := make([]byte, *length)
		for j:=0 ; j<*length ; j++ {
			l := random.Intn(4)
			str[j] = alpha[l]
		}
		seq.SetSequence(str)

		// Write it
		seqOut.Write(*seq)
		seqOut.CheckPanic()
	}
}
