package kmer

// Single word structure
type Word32 struct {
	W uint32
}

// Kmer manager structure
type Kmer32 struct {
	W Word32   // Current word
	K uint8    // K value
	B []uint32 // Base converstion array
}

// Initalize a Kmer manager
func NewKmer32(k int) *Kmer32 {
	var km Kmer32
	km.K = uint8(k)
	km.W = Word32{uint32(0)}
	km.B = make([]uint32, 256)
	
	// No need to treat A => 0
	baseC = uint32(1)
	baseG = uint32(2)
	baseT = uint32(3)
	for i:=1 ; i<k ; i++ {
		baseC = baseC << 2
		baseG = baseG << 2
		baseT = baseT << 2
	}

	// Set base values
	km.B['C'] = baseC
	km.B['c'] = baseC
	km.B['G'] = baseG
	km.B['g'] = baseG
	km.B['T'] = baseT
	km.B['t'] = baseT

	// Return the manager
	return &km
}

// Get Kmer (in bytes)
func (kmer *Kmer32)GetKmer() []byte {
	kmb := make([]byte, kmer.K)
	wrd := kmer.W.W
	for i:=1 ; i<=int(kmer.K) ; i++ {		
		kmb[int(kmer.K)-i] = wrd & uint32(3)
		wrd >> 2
	}
	return kmb
}