package main

import (
	"flag"
	"fmt"
	//"github.com/hdevillers/go-dl-seq/kmer"
)

const (
	MaxStackPerRun int = 250
	MaxStackGroups int = 100
	CounterSize    int = 10
)

type Counter struct {
	data [][]int
	temp [][]int
	coun []int
	/*sem0 chan int
	sem1 chan int
	sem2 chan int
	sem3 chan int
	sem4 chan int
	sem5 chan int*/
}

var sem0 chan int
var sem1 chan int
var sem2 chan int
var sem3 chan int
var sem4 chan int
var sem5 chan int
var sem6 chan int

func NewCounter() *Counter {
	var c Counter
	c.data = make([][]int, MaxStackGroups)
	c.temp = make([][]int, MaxStackGroups)
	c.coun = make([]int, CounterSize)
	/*
		c.sem0 = make(chan int)
		c.sem1 = make(chan int, t) // Limited thread bottle neck
		c.sem2 = make(chan int)
		c.sem3 = make(chan int, 1) // Merging bottle neck
		c.sem4 = make(chan int)    // End of lauching channel
		c.sem5 = make(chan int)    // Final channel
	*/
	return &c
}

func initChanels() {
	for i := 0; i < MaxStackGroups; i++ {
		sem0 <- i
	}
}

// Data loader (With buffer)
func (c *Counter) load(n int) {
	grp := <-sem0
	c.data[grp] = make([]int, MaxStackPerRun)
	sid := 0
	nproc := 0
	for i := 0; i < n; i++ {
		c.data[grp][sid] = i
		sid++
		if sid == MaxStackPerRun {
			// Launch Enumeration
			sem1 <- grp
			go c.enumerate()
			nproc++
			grp = <-sem0
			sid = 0
			c.data[grp] = make([]int, MaxStackPerRun)
		}
	}
	if sid > 0 {
		sem1 <- grp
		go c.enumerate()
		nproc++
	}
	sem4 <- nproc
	c.finish()
}

// First step of the treatment (ex., init temp. counters)
func (c *Counter) enumerate() {
	grp := <-sem1 // Get the group number

	// Init the temp count
	c.temp[grp] = make([]int, CounterSize)

	// Send grp to counter
	sem2 <- grp
	c.count()
}

func (c *Counter) count() {
	grp := <-sem2
	ind := 0
	for i := 0; i < MaxStackPerRun; i++ {
		c.temp[grp][ind] += c.data[grp][i]
		ind++
		if ind == CounterSize {
			ind = 0
		}
	}
	sem3 <- grp
	c.merge()
}

func (c *Counter) merge() {
	grp := <-sem3 // Buffered at 1 to avoid collision!
	for i := 0; i < CounterSize; i++ {
		c.coun[i] += c.temp[grp][i]
	}
	c.temp[grp] = nil
	c.data[grp] = nil
	sem0 <- grp // Free the treated group
	sem5 <- grp // Indicate end of one process
}

func (c *Counter) finish() {
	// Wait that all processes are launched
	nproc := <-sem4

	// Wait for the end of all processes
	for i := 0; i < nproc; i++ {
		<-sem5
	}

	// Close the init. channel
	close(sem0)

	sem6 <- 1
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

	sem0 = make(chan int)
	sem1 = make(chan int, *threads) // Limited thread bottle neck
	sem2 = make(chan int)
	sem3 = make(chan int, 1) // Merging bottle neck
	sem4 = make(chan int)    // End of lauching channel
	sem5 = make(chan int)    // Final channel
	sem6 = make(chan int)

	// Define stack manager
	c := NewCounter()

	initChanels()
	c.load(*stacks)

	<-sem6

}
