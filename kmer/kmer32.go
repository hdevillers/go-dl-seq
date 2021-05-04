package kmer

// Kmer manager structure
type Kmer32 struct {
	W uint32   // Current word
	K uint8    // K value
	B []uint32 // Base converstion array
	A []byte   // Alphabet
	D uint32   // Deleter
}

// Initalize a Kmer manager
func NewKmer32(k int) *Kmer32 {
	var km Kmer32
	km.K = uint8(k)
	km.W = uint32(0)
	km.B = make([]uint32, 256)
	km.A = []byte{'A', 'C', 'G', 'T'}
	km.D = uint32(3) << uint32(2*k)

	// Set base values
	km.B['C'] = uint32(1)
	km.B['c'] = uint32(1)
	km.B['G'] = uint32(2)
	km.B['g'] = uint32(2)
	km.B['T'] = uint32(3)
	km.B['t'] = uint32(3)

	// Return the manager
	return &km
}

// Get Kmer (in bytes)
func (kmer *Kmer32) GetKmer() []byte {
	kmb := make([]byte, kmer.K)
	wrd := kmer.W
	for i := 1; i <= int(kmer.K); i++ {
		kmb[int(kmer.K)-i] = kmer.A[int(wrd&uint32(3))]
		wrd = wrd >> 2
	}
	return kmb
}

// Initialize
func (kmer *Kmer32) Init(b []byte) {
	//Reset
	kmer.W = uint32(0)

	// Push the K first letters
	for i := 0; i < int(kmer.K); i++ {
		kmer.W = (kmer.W << 2) | kmer.B[b[i]]
	}
}

// Add a base in the kmer
func (kmer *Kmer32) AddBase(b byte) {
	// Add the new byte
	kmer.W = (kmer.W << 2) | kmer.B[b]
	// Delete the byte at the left
	kmer.W = (kmer.W | kmer.D) ^ kmer.D
}
