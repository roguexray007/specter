package main

import (
	"fmt"
	"sync"
)

func main() {
	doWork := func(wg *sync.WaitGroup, start <-chan interface{}, id int) {
		go func() {
			defer wg.Done()
			<-start
			fmt.Printf("gorutine no %d completed it's work\n", id)
		}()
	}

	var start chan interface{}
	start = make(chan interface{})

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		doWork(&wg, start, i)
	}
	fmt.Println("unblocking goroutines")
	close(start)
	/*
		closing the channel is both cheaper and faster than performing n writes to signal
		n blocked goroutines. We can use sync.Cond for same behaviour. Infact sync.Cond
		enables repeated signals to unblock multiple goroutines unlike channels.
		check v7

	*/
	wg.Wait() // created a join point
}
