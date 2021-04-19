package fasta

import (
	"errors"
	"strings"

	"github.com/hdevillers/go-dl-seq/seq"
	"github.com/hdevillers/go-dl-seq/seqio/seqitf"
)

const (
	IdPreffix  byte = '>'
	LineLength int  = 60
)

// Fasta sequence reader struct
type Reader struct {
	scan     seqitf.FileScanner
	currId   string
	currDesc string
	eof      bool
}

// Fasta sequence write struct
type Writer struct {
	write seqitf.FileWriter
	Count int
}

// Generate a new reader
func NewReader(fs seqitf.FileScanner) *Reader {
	return &Reader{
		scan:     fs,
		currId:   "",
		currDesc: "",
		eof:      false,
	}
}

// Generate a new writer
//func NewWriter(wf *bufio.Writer) *Writer {
func NewWriter(fw seqitf.FileWriter) *Writer {
	return &Writer{
		write: fw,
		Count: 0,
	}
}

func parseIdLine(idl string) (string, string) {
	data := strings.SplitN(idl, " ", 2)
	if len(data) == 2 {
		return data[0], data[1]
	}
	return data[0], ""
}

// Return true if reachs the end-of-file
func (r *Reader) IsEOF() bool {
	return r.eof
}

// Read a single fasta entry
func (r *Reader) Read() (seq.Seq, error) {
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
					return newSeq, errors.New("[FASTA READER]: Empty sequence or bad format.")
				}

				// Set sequence data
				newSeq.SetId(r.currId)
				newSeq.SetDesc(r.currDesc)

				// Save the new ID
				r.currId, r.currDesc = parseIdLine(string(line[1:]))

				// Return the completed sequence
				return newSeq, nil
			} else {
				// Save the new ID
				r.currId, r.currDesc = parseIdLine(string(line[1:]))

				// Thow an error if the sequence is not nil
				if newSeq.Length() > 0 {
					return newSeq, errors.New("[FASTA READER]: Sequence without ID or possible bad format.")
				}

				// Continue
			}
		} else {
			//TODO: Control input character
			newSeq.AppendSequence(line)
		}
	}
	// Scanning is finicher
	r.eof = true

	// Set last sequence ID and Description
	newSeq.SetId(r.currId)
	newSeq.SetDesc(r.currDesc)

	// Check if the last sequence is empty
	if newSeq.Length() == 0 {
		return newSeq, errors.New("[FASTA READER]: The last sequence is empty.")
	}

	// Return with no error
	return newSeq, nil
}

func (w *Writer) Write(s seq.Seq) error {
	//Add the sequence ID
	if s.Id == "" {
		return errors.New("[FASTA WRITER]: Missing sequence ID.")
	}
	_, err := w.write.Write([]byte(">" + s.Id))
	if err != nil {
		return err
	}

	// Add the description
	if s.Desc != "" {
		_, err = w.write.Write([]byte(" " + s.Desc))
		if err != nil {
			return err
		}
	}
	_, err = w.write.Write([]byte{'\n'})
	if err != nil {
		return err
	}

	// Add the sequence
	// NOTE: We assume that if no error occured above, io.writer is OK
	n := 0
	for i := 0; i < s.Length(); i++ {
		w.write.Write([]byte{s.Sequence[i]})
		n++
		if n == LineLength {
			w.write.Write([]byte{'\n'})
			n = 0
		}
	}
	if n != 0 {
		w.write.Write([]byte{'\n'})
	}

	// Flush written bytes
	err = w.write.Flush()
	w.Count++

	return err
}
