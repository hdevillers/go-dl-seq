package main

import (
	"flag"
	"fmt"
	//"github.com/hdevillers/go-dl-seq/kmer"
)

const (
	MaxStackPerRun int = 250
	MaxStackGroups int = 100
)

type Counter struct {
	data [][]int
	sem1 chan int
	sem2 chan int
	sem3 chan int
}

func NewCounter(t int) *Counter {
	var c Counter
	c.data = make([][]int, MaxStackGroups)
	// Init the first group
	c.data[0] = make([]int, MaxStackPerRun)
	c.sem1 = make(chan int, t)
	c.sem2 = make(chan int, t)
	c.sem3 = make(chan int)
	return &c
}

func (c *Counter) enumerate() {
	grp := <-c.sem1 // Get the group number
	for i := 0; i < MaxStackPerRun; i++ {
		c.data[grp][i] = c.data[grp][i] + 3
	}
	c.sem2 <- grp
	go c.count()
}

func (c *Counter) count() {
	grp := <-c.sem2
	for i := 0; i < MaxStackPerRun; i++ {
		c.data[grp][i] = c.data[grp][i] - 1
	}
	c.sem3 <- grp
}

func main() {
	threads := flag.Int("t", 4, "Number of threads.")
	stacks := flag.Int("s", 1000, "Number of stacks.")
	flag.Parse()

	// Define stack manager
	c := NewCounter(*threads)
	grp := 0 // Current group
	sid := 0 // Current stak

	// Load stack
	for i := 0; i < *stacks; i++ {
		c.data[grp][sid] = i
		sid++
		if sid == MaxStackPerRun {
			// Launch Enumeration
			c.sem1 <- grp
			go c.enumerate()
			grp++
			sid = 0
			c.data[grp] = make([]int, MaxStackPerRun)
		}
	}
	c.sem1 <- grp
	c.enumerate()
	grp++

	for i := 0; i < grp; i++ {
		g := <-c.sem3
		for j := 0; j < MaxStackPerRun; j++ {
			fmt.Println("Data Group\t", g, "\tStack\t", j, "\tValue\t", c.data[g][j])
		}
	}
}
