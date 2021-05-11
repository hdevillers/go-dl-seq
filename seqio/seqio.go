package seqio

import (
	"bufio"
	//"compress/gzip"
	"errors"
	"os"

	gzip "github.com/klauspost/pgzip"

	"github.com/hdevillers/go-dl-seq/seq"
	"github.com/hdevillers/go-dl-seq/seqio/fasta"
	"github.com/hdevillers/go-dl-seq/seqio/fastnq"
	"github.com/hdevillers/go-dl-seq/seqio/fastq"
	"github.com/hdevillers/go-dl-seq/seqio/seqitf"
)

const (
	defaultCompress = false
)

// Reader structure
type Reader struct {
	fcloser seqitf.FileCloser
	sreader seqitf.SeqReader
	seq     seq.Seq
	err     error
}

// Writer structure
type Writer struct {
	fcloser seqitf.FileCloser
	swriter seqitf.SeqWriter
	err     error
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
	var fs seqitf.FileScanner
	var fc seqitf.FileCloser

	if compress[0] {
		// Need de-compression
		fgzip, err := gzip.NewReader(f)
		if err != nil {
			return &Reader{
				err: err,
			}
		}
		fc = fgzip
		fs = bufio.NewScanner(fgzip)
	} else {
		// No de-compression needed
		fs = bufio.NewScanner(f)
		fc = f
	}

	var sreader seqitf.SeqReader
	switch format {
	case "fasta", "fa":
		sreader = fasta.NewReader(fs)
		return &Reader{
			fcloser: fc,
			sreader: sreader}
	case "fastq", "fq":
		sreader = fastq.NewReader(fs)
		return &Reader{
			fcloser: fc,
			sreader: sreader,
		}
	case "fastnq", "fnq":
		sreader = fastnq.NewReader(fs)
		return &Reader{
			fcloser: fc,
			sreader: sreader,
		}
	default:
		return &Reader{
			err: errors.New("[SEQIO READER]: Unsupported format (" + format + ")."),
		}
	}
}

// Read next sequence
func (r *Reader) Next() bool {
	if r.sreader.IsEOF() {
		return false
	} else {
		r.seq, r.err = r.sreader.Read()
		// NOTE: Some parsers return an empty sequence at the end with out error
		if r.err == nil && r.seq.Length() == 0 {
			return false
		} else {
			return true
		}
	}
}

// Get the current sequence
func (r *Reader) Seq() seq.Seq {
	return r.seq
}

// Close file handle
func (r *Reader) Close() {
	r.fcloser.Close()
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
	var fw seqitf.FileWriter
	var fc seqitf.FileCloser

	if compress[0] {
		// Need de-compression
		fgz := gzip.NewWriter(f)
		fw = fgz
		fc = fgz
	} else {
		// No de-compression needed
		fw = bufio.NewWriter(f)
		fc = f
	}

	var swriter seqitf.SeqWriter
	switch format {
	case "fasta", "fa":
		swriter = fasta.NewWriter(fw)
		return &Writer{
			fcloser: fc,
			swriter: swriter,
		}
	case "fastq", "fq", "fastnq", "fnq":
		swriter = fastq.NewWriter(fw)
		return &Writer{
			fcloser: fc,
			swriter: swriter,
		}
	default:
		return &Writer{
			err: errors.New("[SEQIO WRITER]: Unsupported format (" + format + ")."),
		}
	}
}

// Append a sequence in the output file
func (w *Writer) Write(s seq.Seq) {
	w.err = w.swriter.Write(s)
}

// Close output file
func (w *Writer) Close() {
	err := w.swriter.Flush()
	w.err = err
	w.fcloser.Close()
}

// Throw a panic in case of error
func (w *Writer) CheckPanic() {
	if w.err != nil {
		panic(w.err)
	}
}
