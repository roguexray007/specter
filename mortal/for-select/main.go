package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {

	/*
		Notes on channel
		1. closing a channel to unblock multiple goroutines at the same time. But reproducing the behaviour
		of repeated unblock would be more difficult via channel. Sync.cond is better here.
	*/

	// Pattern 1. sending iteration variable out on a channel

	/*
			Let's take a look at how we organize different types of channels to begin buidling
			something that's robust and stable.
			The first thing we should do to put channels in the right context is to assign channel
			ownership.I'll define ownership as being a goroutine that instantiates , writes, and closes
			a channel. It's important to clarify which goroutines own a channel in order to reason about
			our programs logically. (Unidirectional channel declarations are the tool that will allow us
			to distinguish between goroutines that own channels and those that only utilize them.

			The goroutine that owns a channel should
			1. instantiate the channel
			2. perform writes or pass ownership to another goroutine
			3. close the channel

			Let's look at blocking operations that can occur when reading. As a consumer of channel,
			I have to worry about 2 things
			1. Knowing when a channel is closed
			2. Responsibly handling blocking for any reason

		by assigning the responsibilities to channel owners, a few things happen
		1. Because we are the one initializing the channel, we remove the risk of deadlocking by
			writing to a nil channel. (write/read to nil channel is blocking)
		2. Because we are the one initializing the channel, we remove the risk of panicing by closing
			a nil channel.
		3. Because we are the one who decides when the channel gets closed, we remove the risk of
			panicing by writing to a closed channel
		4. Because we are the one who decides when the channel gets closed, we remove the risk of
			panicing by closing a channel more than once
		5. We wield the type checker at compile time to prevent improper writes to our channel

		lexical confinement - exposing read or write aspects of a channel to the concurrent
		processes that need them

		Here we start an anonymous goroutine that performs writes on output, Notice that how we've inverted
		how we create goroutines. It is now encapsulated within the surrounding function
	*/
	chanOwner := func(done <-chan interface{}) <-chan string {

		/*
			Separate declaration and instantiation
			this is somewhat interesting because it means that the goroutine that instantiates a channel
			controls whether it's buffered. This suggests the creation of a channel should probably
			be tightly coupled to goroutines that will be performing writes on it so that we can reason
			about its behaviour and performance more easily

			It also bears mentioning that if a buffered channel is empty and has a receiver, the buffer will
			be bypassed and the value will be directly passed from the sender to the receiver
		*/

		var output chan string
		output = make(chan string)
		go func() {
			defer close(output)
			defer func() { fmt.Println("exiting after sending entire data") }()
			for _, val := range []string{"hi", "specter", "welcome", "to", "party"} {
				select {
				case output <- val:
				case <-done:
					return
				}
			}
		}()
		return output
	}

	consumer := func(output <-chan string) {
		for val := range output {
			fmt.Println(val)
		}
	}

	output := chanOwner(nil)
	consumer(output)

	fmt.Println("pattern2 ahead")

	// 2. Looping infinitely waiting to be stopped
	pattern2 := func(wg *sync.WaitGroup, done <-chan interface{}) {
		defer wg.Done()
		for {
			select {
			case <-done:
				return
			default:
				time.Sleep(200 * time.Millisecond)
				fmt.Println("pattern2 in work")
			}
		}
	}

	doneOwner := func(wg *sync.WaitGroup) <-chan interface{} {
		var done chan interface{}
		done = make(chan interface{})
		go func() {
			defer wg.Done()
			select {
			case <-time.After(1 * time.Second):
				fmt.Println("closing done")
				close(done)
			}
		}()

		return done
	}

	var wg sync.WaitGroup
	wg.Add(2)
	done := doneOwner(&wg)
	go pattern2(&wg, done)
	wg.Wait()
}
