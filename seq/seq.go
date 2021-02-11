package seq

type Seq struct {
	id       string
	sequence []byte
}

func NewSeq(id string, sequence []byte) *Seq {
	p := Seq{id: id, sequence: sequence}
	return &p
}

func (s *Seq) Length() int {
	return len(s.sequence)
}
