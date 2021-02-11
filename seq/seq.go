package seq

type Seq struct {
	id       string
	desc     string
	sequence []byte
}

func NewSeq(id string, sequence []byte) *Seq {
	p := Seq{id: id, sequence: sequence}
	return &p
}

func (s *Seq) SetId(id string) {
	s.id = id
}

func (s *Seq) SetDesc(desc string) {
	s.desc = desc
}

func (s *Seq) SetSequence(sequence []byte) {
	s.sequence = sequence
}

func (s *Seq) Length() int {
	return len(s.sequence)
}
