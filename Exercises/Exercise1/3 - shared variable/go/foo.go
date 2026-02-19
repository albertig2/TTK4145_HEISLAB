// Use `go run foo.go` to run your program

package main

import (
	. "fmt"
	"runtime"
	"time"
)

var i = 0

var inc = make(chan struct{})
var dec = make(chan struct{})
var done = make(chan struct{})
var get = make(chan chan int)

func i_server() {
	for {
		select {
		case <-inc:
			i++
		case <-dec:
			i--
		case reply := <-get:
			reply <- i
		}
	}
}

func incrementing() {
	for k := 0; k < 1000000; k++ {
		inc <- struct{}{}
	}
	done <- struct{}{}
	//TODO: increment i 1000000 times
}

func decrementing() {
	for k := 0; k < 1000000; k++ {
		dec <- struct{}{}
	}
	done <- struct{}{}
	//TODO: decrement i 1000000 times
}

func main() {
	// What does GOMAXPROCS do? What happens if you set it to 1?
	runtime.GOMAXPROCS(3)
	// TODO: Spawn both functions as goroutines
	go i_server()
	go incrementing()
	go decrementing()

	// We have no direct way to wait for the completion of a goroutine (without additional synchronization of some sort)
	// We will do it properly with channels soon. For now: Sleep.

	<-done
	<-done

	reply := make(chan int)
	get <- reply
	i := <-reply

	time.Sleep(500 * time.Millisecond)
	Println("The magic number is:", i)
}
