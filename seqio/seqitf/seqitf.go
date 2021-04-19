package seqitf

import (
	"github.com/hdevillers/go-dl-seq/seq"
)

/*
	This provide a set of interfaces requires in seqio module and
	in the different sequence format parsers.
*/

// Generic interface to read sequences
type SeqReader interface {
	Read() (seq.Seq, error)
	IsEOF() bool
}

// Generic interface to write sequences
type SeqWriter interface {
	Write(seq.Seq) error
}

// File handling interfaces
type FileWriter interface {
	Write(p []byte) (n int, err error)
	Flush() error
}
type FileCloser interface {
	Close() error
}
