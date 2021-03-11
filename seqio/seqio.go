package seqio

import (
	"errors"
	"os"

	"github.com/hdevillers/go-dl-seq/seq"
	"github.com/hdevillers/go-dl-seq/seqio/fasta"
)

// Reader structure
type Reader struct {
	fhdl  *os.File
	fread interface {
		Read() (seq.Seq, error)
		IsEOF() bool
	}
	seq seq.Seq
	err error
}

// Writer structure
type Writer struct {
	fhdl  *os.File
	fwrit interface {
		Write(seq.Seq) error
	}
	err error
}

// Create a new reader (from a file name and a format)
func NewReader(file string, format string) *Reader {
	// Open file in read mode
	f, err := os.Open(file)
	if err == nil {
		switch format {
		case "fasta", "fa":
			var fread fasta.ReaderInterface
			fread = fasta.NewReader(f)
			return &Reader{
				fhdl:  f,
				fread: fread,
			}
		default:
			return &Reader{
				err: errors.New("[SEQIO READER]: Unsupported format (" + format + ")."),
			}
		}
	} else {
		return &Reader{
			err: err,
		}
	}
}

// Read next sequence
func (r *Reader) Next() bool {
	if r.fread.IsEOF() {
		return false
	} else {
		r.seq, r.err = r.fread.Read()
		return true
	}
}

// Get the current sequence
func (r *Reader) Seq() seq.Seq {
	return r.seq
}

// Close file handle
func (r *Reader) Close() {
	r.fhdl.Close()
}

// Get errors
func (r *Reader) CheckPanic() {
	if r.err != nil {
		panic(r.err)
	}
}

// Create a new Writer (from a file name and a format)
func NewWriter(file string, format string) *Writer {
	// Open a file in write/overide mode
	var f *os.File
	var err error
	// Write into stdout if file is empty
	if file == "" {
		f = os.Stdout
	} else {
		f, err = os.Create(file)
		if err != nil {
			return &Writer{
				err: err,
			}
		}
	}
	switch format {
	case "fasta","fa":
		var fwrit fasta.WriterInterface
		fwrit = fasta.NewWriter(f)
		return &Writer{
			fhdl:  f,
			fwrit: fwrit,
		}
	default:
		return &Writer{
			err: errors.New("[SEQIO WRITER]: Unsupported format (" + format + ")."),
		}
	}
}

// Append a sequence in the output file
func (w *Writer) Write(s seq.Seq) {
	w.err = w.fwrit.Write(s)
}

// Close output file
func (w *Writer) Close() {
	w.fhdl.Close()
}

// Throw a panic in case of error
func (w *Writer) CheckPanic() {
	if w.err != nil {
		panic(w.err)
	}
}