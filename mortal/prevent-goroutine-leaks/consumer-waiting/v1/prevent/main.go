package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	/*
		goroutines are not garbage collected by go runtime, so regardless of how small their memory
		footprint is, we don't want to leave them lying about our process. So, how do we go about
		ensuring they are cleaned up?

		Why would a goroutine exit?
		1. when it completed it's work.
		2. when it cannot continue to work due to an unrecoverable error.
		3. when it's told to stop working.

		first two are your algorithm. What about work cancellation.

	*/

	/*
		Here we see that main goroutine passes a nil channel to consumer. Therefore, the input channel
		will never actually gets any strings written onto it, and the consumer goroutine will continue
		remain in memory for the lifetime of this process. (we would even deadlock if we joined
		the main goroutine and consumer goroutine.

		The way to successfully mitigate this is to establish a signal b/w the parent gorutine and child
		goroutine that allows parent to signal cancellation to its children. By convention the signal is
		usually a read-only channel named done
	*/
	consumer := func(wg *sync.WaitGroup, input <-chan string, done <-chan interface{}) {
		go func() {
			defer wg.Done()
			defer fmt.Println("closing consumer")
			for {
				select {
				case <-done:
					return
				case val := <-input:
					fmt.Println("received :", val)
				}
			}
		}()
	}

	var wg sync.WaitGroup
	doneOwner := func(wg *sync.WaitGroup) <-chan interface{} {
		var done chan interface{}
		done = make(chan interface{})
		go func() {
			defer func() {
				close(done)
				fmt.Println("closed done channel")
				wg.Done()
			}()
			select {
			case <-time.After(time.Second * 2):
				return
			}
		}()
		return done
	}

	wg.Add(1)
	done := doneOwner(&wg)

	wg.Add(1)
	consumer(&wg, nil, done)

	wg.Wait() // no deadlock even when nil channel is sent as consumer goroutine is sent a signal to shutdown
	// after 2 sec by doneOwner goroutine
	fmt.Println("exiting ")

}
