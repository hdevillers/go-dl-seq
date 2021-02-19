package encoding

import (
	"fmt"
	"errors"
)

type OneHot struct {
	seq []byte
	oho []uint16
	trs map[byte]int
	err error
}

func initTrs() map[byte]int {
	return map[byte]int{
		'A':0,
		'C':1,
		'G':2,
		'T':3,
	}
}

func NewOneHot(s []byte) *OneHot {
	return &OneHot{
		seq: s,
		trs: initTrs(),
	}
}

func (o *OneHot) SetSeq(s []byte) {
	o.seq = s
}

func (o *OneHot) Compute() {
	if len(o.seq) == 0 {
		o.err = errors.New("[ONEHOT]: No sequence provided.")
	} else {
		oho := make([]uint16, len(o.seq) * 4)
		for i,l := range o.seq {
			oho[o.trs[l]*len(o.seq) + i] = uint16(1)
		}
		o.oho = oho
	}
}

func (o *OneHot) Show() {
	if len(o.seq) == 0 {
		o.err = errors.New("[ONEHOT]: No sequence provided.")
	} else if len(o.oho) == 0 {
		o.err = errors.New("[ONEHOT]: Launch compute method first.")
	} else {
		for i := 0 ; i<4 ; i++ {
			for j := 0 ; j<len(o.seq) ; j++ {
				fmt.Printf("%d", o.oho[i*len(o.seq)+j])
			}
			fmt.Printf("\n")
		}
		fmt.Printf("\n\n")
	}
}

func (o *OneHot) CheckPanic() {
	if o.err != nil {
		panic(o.err)
	}
}

func (o *OneHot) GetOho() []uint16 {
	return o.oho
}