package main

import (
	"flag"

	"github.com/hdevillers/go-dl-seq/seqio"
)

func main() {
	input := flag.String("i", "", "Input fastq file.")
	output := flag.String("o", "", "Ouput fastq file.")
	minL := flag.Int("min-length", 10, "Mininal read length.")
	maxL := flag.Int("max-length", -1, "Maximal read length.")
	maxN := flag.Int("max-reads", -1, "Maximal number of kept reads.")
	skipN := flag.Bool("skip-n", false, "Skip reads with N.")
	d := flag.Bool("d", false, "Decompress input.")
	c := flag.Bool("c", false, "Compress output.")
	flag.Parse()

	if *input == "" {
		panic("You must provide an input fastq file.")
	}

	seqIn := seqio.NewReader(*input, "fastq", *d)
	seqOut := seqio.NewWriter(*output, "fastq", *c)
	defer seqIn.Close()
	defer seqOut.Close()
	n := 0

	for seqIn.Next() {
		seqIn.CheckPanic()
		s := seqIn.Seq()

		// Check min length
		if s.Length() < *minL {
			continue
		}

		// Check max length
		if *maxL > 0 && s.Length() > *maxL {
			continue
		}

		// Check for N
		if *skipN {
			for _, b := range s.Sequence {
				if b == 'N' || b == 'n' {
					continue
				}
			}
		}

		// All checks passed
		seqOut.Write(s)
		n++
		if *maxN > 0 {
			if n == *maxN {
				break
			}
		}
	}

}
