// Use `go run foo.go` to run your program

package main

import (
	. "fmt"
	"runtime"
)

var i = 0

func incrementing(ch1 chan int, ch2 chan int) {
	//TODO: increment i 1000000 times
	for j := 0; j < 1000001; j++ {
		ch1 <- 1
	}
	ch2 <- 1
}

func decrementing(ch1 chan int, ch2 chan int) {
	//TODO: decrement i 1000000 times
	for j := 0; j < 1000000; j++ {
		ch1 <- 1
	}
	ch2 <- 1
}

func synchronizer(ch1 chan int, ch2 chan int, ch_read chan int) {
	for {
		select {
		case <-ch1:
			i++
		case <-ch2:
			i--
		case ch_read <- i:
		}
	}
}

func main() {
	// create channels
	ch_data1 := make(chan int)
	ch_data2 := make(chan int)
	ch_ack := make(chan int)
	ch_read := make(chan int)

	// What does GOMAXPROCS do? What happens if you set it to 1?
	runtime.GOMAXPROCS(3)

	// TODO: Spawn both functions as goroutines
	go synchronizer(ch_data1, ch_data2, ch_read)
	go incrementing(ch_data1, ch_ack)
	go decrementing(ch_data2, ch_ack)

	select {
	case <-ch_ack:
		// do nothing
	}

	select {
	case <-ch_ack:
		// do nothing
	}

	// We have no direct way to wait for the completion of a goroutine (without additional synchronization of some sort)
	// We will do it properly with channels soon. For now: Sleep.
	// time.Sleep(500 * time.Millisecond)
	Println("The magic number is:", <-ch_read)
}
