package fastnq

import (
	"github.com/hdevillers/go-dl-seq/seq"
	"github.com/hdevillers/go-dl-seq/seqio/seqitf"
)

const (
	IdPreffix byte = '@'
	SpPreffix byte = '+'
)

// Fastq sequence reader struct
type Reader struct {
	scan     seqitf.FileScanner
	currId   string
	eof      bool
	waitQual bool
}

// Generate a new reader
func NewReader(fs seqitf.FileScanner) *Reader {
	return &Reader{
		scan:     fs,
		currId:   "",
		eof:      false,
		waitQual: false,
	}
}

// Return true if reachs the end-of-file
func (r *Reader) IsEOF() bool {
	return r.eof
}

// Read a single fastq entry
func (r *Reader) Read() (seq.Seq, error) {
	// Initialize the new sequence
	var newSeq seq.Seq

	/*
		NOTE: this parser version is made to save time:
		1) Quality is not considered
		2) Fastq file is supposed well formed
	*/

	for r.scan.Scan() {
		// Check possible scanning error
		err := r.scan.Err()
		if err != nil {
			return newSeq, err
		}

		// Get the ID line
		line := r.scan.Bytes()
		newSeq.SetId(string(line[1:]))

		// Get the sequence line
		r.scan.Scan()
		line = r.scan.Bytes()
		newSeq.AppendSequence(line)

		// Skip spacer line and quality line
		r.scan.Scan()
		r.scan.Scan()

		return newSeq, nil
	}

	r.eof = true
	return newSeq, nil
}
