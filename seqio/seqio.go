package seqio

import (
	"bufio"
	"compress/gzip"
	"errors"
	"os"

	"github.com/hdevillers/go-dl-seq/seq"
	"github.com/hdevillers/go-dl-seq/seqio/fasta"
	"github.com/hdevillers/go-dl-seq/seqio/fastq"
)

const (
	defaultCompress = false
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
func NewReader(file string, format string, compress ...bool) *Reader {
	// Open file in read mode
	var f *os.File
	var err error
	if file == "STDIN" {
		f = os.Stdin
	} else {
		f, err = os.Open(file)
		if err != nil {
			return &Reader{
				err: err,
			}
		}
	}

	// Check compression argument
	if len(compress) == 0 {
		compress = append(compress, defaultCompress)
	}

	// Inti. the bufio.Scanner
	var sf *bufio.Scanner
	if compress[0] {
		// Need de-compression
		rf, err := gzip.NewReader(f)
		if err != nil {
			return &Reader{
				err: err,
			}
		}
		sf = bufio.NewScanner(rf)
	} else {
		// No de-compression needed
		sf = bufio.NewScanner(f)
	}

	switch format {
	case "fasta", "fa":
		var fread fasta.ReaderInterface
		fread = fasta.NewReader(sf)
		return &Reader{
			fhdl:  f,
			fread: fread}
	case "fastq", "fq":
		var fread fastq.ReaderInterface
		fread = fastq.NewReader(sf)
		return &Reader{
			fhdl:  f,
			fread: fread,
		}
	default:
		return &Reader{
			err: errors.New("[SEQIO READER]: Unsupported format (" + format + ")."),
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
func NewWriter(file string, format string, compress ...bool) *Writer {
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

	// Check compression argument
	if len(compress) == 0 {
		compress = append(compress, defaultCompress)
	}

	// Inti. the bufio.Scanner
	var wf *bufio.Writer
	if compress[0] {
		// Need de-compression
		rf := gzip.NewWriter(f)
		wf = bufio.NewWriter(rf)
	} else {
		// No de-compression needed
		wf = bufio.NewWriter(f)
	}

	switch format {
	case "fasta", "fa":
		var fwrit fasta.WriterInterface
		fwrit = fasta.NewWriter(wf)
		return &Writer{
			fhdl:  f,
			fwrit: fwrit,
		}
	case "fastq", "fq":
		var fwrit fastq.WriterInterface
		fwrit = fastq.NewWriter(wf)
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
