package pattern

import (
	"math/rand"
	"regexp"
	"strconv"
)

type Nucl struct {
	alpha   [][]byte
	iupac   map[byte]int
	alen    []int
	plen    int
	pattern []int 
}

func initNuclAlpha() [][]byte {
	return [][]byte{
		{'A', 'C', 'G', 'T'},
		{'A'},
		{'C'},
		{'G'},
		{'T'},
		{'A', 'G'},
		{'C', 'T'},
		{'G', 'C'},
		{'A', 'T'},
		{'G', 'T'},
		{'A', 'C'},
		{'C', 'G', 'T'},
		{'A', 'G', 'T'},
		{'A', 'C', 'T'},
		{'A', 'C', 'G'},
	}
}

func initNuclIupac() map[byte]int {
	return map[byte]int{
		'N': 0,
		'A': 1,
		'C': 2,
		'G': 3,
		'T': 4,
		'R': 5,
		'Y': 6,
		'S': 7,
		'W': 8,
		'K': 9,
		'M': 10,
		'B': 11,
		'D': 12,
		'H': 13,
		'V': 14,
	}
}

func NewNucl(l int) *Nucl {
	n := Nucl{
		alpha:   initNuclAlpha(),
		iupac:   initNuclIupac(),
		plen:    l,
		pattern: make([]int, l),
	}
	// Initialize pattern with 0 (=> N)
	for i:=0 ; i<l ; i++ {
		n.pattern[i] = 0
	}
	// Initialize alen (length of each sub alpha)
	n.alen = make([]int, len(n.alpha))
	for i:=0 ; i<len(n.alpha) ; i++ {
		n.alen[i] = len(n.alpha[i])
	}
	return &n
}

func (n *Nucl)EditNuclPattern(s string) error {
	// Copy pattern
	p := n.pattern

	// Parse the input string
	re := regexp.MustCompile(`([0-9]+)([A-Z]+)`)
	in := re.FindAllStringSubmatch(s, -1)

	// Edit pattern
	for i:=0 ; i<len(in) ; i++ {
		// Convert string to byte
		pa := []byte(in[i][2])
		at, err := strconv.Atoi(in[i][1])
		if err != nil {
			return err
		}
		at--
		for j:=0 ; j<len(pa) && at<n.plen ; j++ {
			p[at] = n.iupac[pa[j]]
			at++
		}
	}
	n.pattern = p
	return nil
}

func (n *Nucl)RandomNucl(r *rand.Rand) []byte {
	// Initialze the random sequence
	rs := make([]byte, n.plen)

	// Select a random nucleotide according to the pattern
	for i:=0 ; i<n.plen ; i++ {
		rs[i] = n.alpha[n.pattern[i]][r.Intn(n.alen[n.pattern[i]])]
	}

	// Return the new sequence
	return rs
}