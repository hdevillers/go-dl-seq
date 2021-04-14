package seq

type Seq struct {
	Id       string
	Desc     string
	Sequence []byte
	Quality  []byte
}

func NewSeq(id string) *Seq {
	p := Seq{Id: id}
	return &p
}

func (s *Seq) SetId(id string) {
	s.Id = id
}

func (s *Seq) SetDesc(desc string) {
	s.Desc = desc
}

func (s *Seq) SetSequence(sequence []byte) {
	s.Sequence = sequence
}

func (s *Seq) SetQuality(quality []byte) {
	s.Quality = quality
}

func (s *Seq) AppendSequence(sequence []byte) {
	s.Sequence = append(s.Sequence, sequence...)
}

func (s *Seq) AppendQuality(quality []byte) {
	s.Quality = append(s.Quality, quality...)
}

func (s *Seq) Length() int {
	return len(s.Sequence)
}
