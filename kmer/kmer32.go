package kmer

// Kmer manager structure
type Kmer32 struct {
	W uint32   // Current word
	K uint8    // K value
	B []uint32 // Base converstion array
	A []byte   // Alphabet
	D uint32   // Deleter
}

// Initialize alphabet
func initAlphabet() []byte {
	return []byte{'A', 'C', 'G', 'T'}
}

// Initalize a Kmer manager
func NewKmer32(k int) *Kmer32 {
	var km Kmer32
	km.K = uint8(k)
	km.W = uint32(0)
	km.B = make([]uint32, 256)
	km.A = initAlphabet()
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

func Kmer32String(w uint32, k int) string {
	// Get alphabet
	a := initAlphabet()
	ws := make([]byte, k)
	for i := 1; i <= k; i++ {
		ws[k-i] = a[int(w&uint32(3))]
		w = w >> 2
	}
	return string(ws)
}

// Kmer Slice for uint32
type Kmer32Slice struct {
	W []uint32
}

func NewKmer32Slice() *Kmer32Slice {
	return &Kmer32Slice{}
}

func (ks *Kmer32Slice) Len() int {
	return len(ks.W)
}

func (ks *Kmer32Slice) Less(i, j int) bool {
	return ks.W[i] < ks.W[j]
}

func (ks *Kmer32Slice) Equal(i, j int) bool {
	return ks.W[i] == ks.W[j]
}

func (ks *Kmer32Slice) Swap(i, j int) {
	ks.W[i], ks.W[j] = ks.W[j], ks.W[i]
}

func (ks *Kmer32Slice) Extend(n int) {
	ks.W = append(ks.W, make([]uint32, n)...)
}
