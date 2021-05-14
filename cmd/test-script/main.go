package main

import (
	"flag"
	"fmt"
	//"github.com/hdevillers/go-dl-seq/kmer"
)

const (
	CounterSize int = 10
)

type Counter struct {
	data   []int
	from   []int
	to     []int
	thread int
	temp   [][]int
	coun   []int
}

func NewCounter(t int) *Counter {
	var c Counter
	c.thread = t
	c.temp = make([][]int, t)
	c.from = make([]int, t)
	c.to = make([]int, t)
	c.coun = make([]int, CounterSize)
	return &c
}

// Data loader
func (c *Counter) load(loaded chan int, n int) {
	c.data = make([]int, n)
	for i := 0; i < n; i++ {
		c.data[i] = i
	}

	if c.thread == 1 {
		c.from[0] = 0
		c.to[0] = n - 1
		loaded <- 0
	} else {
		div := n / c.thread
		mod := n % c.thread

		// Set the first coordinates
		c.from[0] = 0
		c.to[0] = div + mod - 1
		loaded <- 0

		// Set the other coordinates
		for i := 1; i < c.thread; i++ {
			c.from[i] = i*div + mod
			c.to[i] = (i+1)*div + mod - 1
			loaded <- i
		}
	}
}

// First step of the treatment (ex., init temp. counters)
func (c *Counter) enumerate(loaded chan int, enumerated chan int) {
	for i := range loaded {
		// Init the temp count
		c.temp[i] = make([]int, CounterSize)

		// Send grp to counter
		enumerated <- i
	}
}

func (c *Counter) count(enumerated chan int, counted chan int) {
	for i := range enumerated {
		ind := 0
		for j := c.from[i]; j <= c.to[i]; j++ {
			c.temp[i][ind] += c.data[j]
			ind++
			if ind == CounterSize {
				ind = 0
			}
		}
		counted <- i
	}
}

func (c *Counter) merge(counted chan int, merged chan int) {
	for i := range counted {
		for j := 0; j < CounterSize; j++ {
			c.coun[j] += c.temp[i][j]
		}

		merged <- i
	}
}

func (c *Counter) print() {
	for i := 0; i < CounterSize; i++ {
		fmt.Println("Value ", i, " is ", c.coun[i])
	}
}

func main() {
	threads := flag.Int("t", 4, "Number of threads.")
	stacks := flag.Int("s", 1000, "Number of stacks.")
	flag.Parse()

	// Init. channeld
	loaded := make(chan int) // Monitor loaded data
	enumerated := make(chan int)
	counted := make(chan int, 1)
	merged := make(chan int)

	c := NewCounter(*threads)

	// Step 1 load data
	go c.load(loaded, *stacks)

	// Step 2 enumerate
	go c.enumerate(loaded, enumerated)

	// Step 3 count
	go c.count(enumerated, counted)

	// Step 4 merge
	go c.merge(counted, merged)

	// Wait all merged counts
	for i := 0; i < *threads; i++ {
		<-merged
	}

	c.print()
}
