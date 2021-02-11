package fasta

import (
	"bufio"
	"errors"
	"io"
	"strings"

	"github.com/hdevillers/go-dl-seq/seq"
)

const (
	IdPreffix = '>'
)

// Fasta sequence reader struct
type Reader struct {
	scan     *bufio.Scanner
	currId   string
	currDesc string
}

// Generate a new reader
func NewReader(f io.Reader) *Reader {
	return &Reader{
		scan:   bufio.NewScanner(f),
		currId: "",
	}
}

func parseIdLine(idl string) (string, string) {
	data := strings.SplitN(idl, " ", 2)
	if len(data) == 2 {
		return data[0], data[1]
	}
	return data[0], ""
}

// Read a single fasta entry
func (r *Reader) Read() (seq.Seq, error) {
	// Initialize the new sequence
	var currSeq []byte
	var newSeq seq.Seq

	for r.scan.Scan() {
		// Check possible scanning error
		err := r.scan.Err()
		if err != nil {
			return newSeq, err
		}

		// Get the scanned line
		line := r.scan.Bytes()
		//fmt.Println(string(line))

		
		if line[0] == IdPreffix {
			// This is an ID line
			if r.currId != "" {
				// Return the current sequence if not nil
				if len(currSeq) == 0 {
					// Empty sequence or bad format
					return newSeq, errors.New("[FASTA]: Empty sequence or bad format.")
				}

				// Set sequence data
				newSeq.SetId(r.currId)
				newSeq.SetDesc(r.currDesc)
				newSeq.SetSequence(currSeq)

				// Save the new ID
				r.currId, r.currDesc = parseIdLine( string(line) )


				return newSeq, nil
			}

		} else {

		}
	}
	return newSeq, nil
}