package kmer

type Kmer32 struct {
	K    int
	Base []byte
}

func NewKmer32(k int) *Kmer32 {
	var km Kmer32
	km.K = k
	km.Base = []byte{'A', 'C', 'G', 'T'}
	return &km
}

func (km *Kmer32) Kmer32ToBytes(w uint32) []byte {
	wb := make([]byte, km.K)
	se := uint32(3) // 2 bits selectors
	for i := 1; i <= km.K; i++ {
		wb[km.K-i] = km.Base[int(w&se)]
		w = w >> 2
	}
	return wb
}
