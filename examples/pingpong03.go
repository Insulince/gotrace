package main

import "time"

func main() {
	var Ball int
	table := make(chan int)

	for i := 0; i < 100; i++ {
		// this sleep is added to more or less preserve an order
		// of the goroutines id on creation
		time.Sleep(10 * time.Millisecond)
		go player(table)
	}

	table <- Ball
	time.Sleep(1 * time.Second)
	<-table
}

func player(table chan int) {
	for {
		ball := <-table
		ball++
		time.Sleep(100 * time.Millisecond)
		table <- ball
	}
}