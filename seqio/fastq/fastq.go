package fastq

import (
	"bufio"
	"errors"
	"io"
	//"strings"

	"github.com/hdevillers/go-dl-seq/seq"
)

const (
	IdPreffix = '@'
	SpPreffix = '+'
)

// Define the fastq reader interface
type ReaderInterface interface {
	Read() (seq.Seq, error)
	IsEOF() bool
}

// Define the fastq writer interface
type WriterInterface interface {
	Write(seq.Seq) error
}

// Fastq sequence reader struct
type Reader struct {
	scan     *bufio.Scanner
	currId   string
	eof      bool
	waitQual bool
}

// Fastq sequence writer struct
type Writer struct {
	write *bufio.Writer
	Count int
}

// Generate a new reader
func NewReader(f io.Reader) *Reader {
	return &Reader{
		scan:     bufio.NewScanner(f),
		currId:   "",
		eof:      false,
		waitQual: false,
	}
}

// Generate a new writer
func NewWriter(f io.Writer) *Writer {
	return &Writer{
		write: bufio.NewWriter(f),
		Count: 0,
	}
}

// Return true if reachs the end-of-file
func (r *Reader) IsEOF() bool {
	return r.eof
}

// Read a single fastq entry
func (r* Reader) Read() (seq.Seq, error) {
	// Initialize the new sequence
	var newSeq seq.Seq

	for r.scan.Scan() {
		// Check possible scanning error
		err := r.scan.Err()
		if err != nil {
			return newSeq, err
		}

		// Get the scanned line
		line := r.scan.Bytes()

		if line[0] == IdPreffix {
			// This is an ID line
			if r.currId != "" {
				// Return the current sequence if not nil
				if newSeq.Length() == 0 {
					// Empty sequence or bad format
					return newSeq, errors.New("[FASTQ READER]: Empty sequence or bad format.")
				}

				// Seq sequence data
				newSeq.SetId(r.currId)

				// Save the new ID
				r.currId = string(line[1:]) // Only skip the line preffix
				r.waitQual = false

				// Return the sequence
				return newSeq, nil
			} else {
				// Save the new ID
				r.currId = string(line[1:])

				// The sequence object should be empty
				if newSeq.Length() > 0 {
					return newSeq, errors.New("[FASTQ READER]: Sequence without ID ou bad format.")
				}

				// Continue
			}
		} else {
			if line[0] == SpPreffix {
				// Finished to read sequence line(s)
				// Start reading the quality
				r.waitQual = true

				// At that step, newSeq.Length must not be null
				if newSeq.Length() == 0 {
					return newSeq, errors.New("[FASTQ READER]: Empty sequence or bad format.")
				}

				// Continue
			} else {
				// Read sequence data or quality data
				// NOTE: We accept non standard fastq with sequence on multiple lines
				if r.waitQual {
					newSeq.AppendQuality(line)
				} else {
					newSeq.AppendSequence(line)
				}
			}
		}
	}
	// Scanning is finished
	r.eof = true

	// Set last sequence ID
	newSeq.SetId(r.currId)

	if newSeq.Length() == 0 {
		return newSeq, errors.New("[FASTQ READER]: Last sequence is null.")
	}

	// Return with no error
	return newSeq, nil
}

func (w *Writer) Write(s seq.Seq) error {
	// Check sequence validity
	if s.Id == "" {
		return errors.New("[FASTQ WRITER]: Missing sequence ID.")
	}
	if s.Length() == 0 {
		return errors.New("[FASTQ WRITER]: Cannot write out empty sequences.")
	}
	if s.Length() != len(s.Quality) {
		return errors.New("[FASTQ WRITER]: Sequence and quality with different lengths.")
	}

	// Add the ID
	_, err := w.write.WriteString(string(IdPreffix) + s.Id)
	if err != nil {
		return err
	}
	err = w.write.WriteByte('\n')
	if err != nil {
		return err
	}

	// Add the sequence
	_, err = w.write.Write(s.Sequence)
	if err != nil {
		return err
	}
	err = w.write.WriteByte('\n')
	if err != nil {
		return err
	}
	err = w.write.WriteByte(SpPreffix)
	if err != nil {
		return err
	}
	err = w.write.WriteByte('\n')
	if err != nil {
		return err
	}

	// Add the quaity
	_, err = w.write.Write(s.Quality)
	if err != nil {
		return err
	}
	err = w.write.WriteByte('\n')
	if err != nil {
		return err
	}

	err = w.write.Flush()
	w.Count++

	return err
}