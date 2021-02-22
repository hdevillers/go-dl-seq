package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"

	"github.com/hdevillers/go-dl-seq/pattern"
	"github.com/hdevillers/go-dl-seq/seq"
	"github.com/hdevillers/go-dl-seq/seqio"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	// Retrieve argument values
	output := flag.String("output", "", "Output file name/path.")
	format := flag.String("format", "fasta", "Output format.")
	length := flag.Int("length", 200, "Required sequence length.")
	count := flag.Int("n", 1, "Number of required sequence(s).")
	base := flag.String("base", "RandSeq_", "Sequence ID base name.")
	seed := flag.Int64("seed", 0, "Random seed initializer.")
	pa := flag.String("pattern", "", "Set specific pattern(s).")
	desc := flag.String("desc", "", "Set a description for each sequence.")
	flag.Parse()

	if *output == "" {
		panic("You must provide an output file name.")
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

	// Create a new pattern generator
	patt := pattern.NewNucl(*length)

	// Edit pattern if required
	if *pa != "" {
		err := patt.EditNuclPattern(*pa)
		check(err)
	}

	// Generate the required sequences
	for i := 0; i < *count; i++ {
		// Create the new ID
		id := *base + fmt.Sprintf("%06d", i)

		// Create the new seq object
		seq := seq.NewSeq(id)

		// Add a sequence
		str := patt.RandomNucl(random)
		seq.SetSequence(str)

		// Add the description if necessary
		seq.SetDesc(*desc)

		// Write it
		seqOut.Write(*seq)
		seqOut.CheckPanic()
	}
}
