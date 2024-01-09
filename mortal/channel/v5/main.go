package main

import "fmt"

func main() {
	var stream chan string
	stream = make(chan string)

	go func() {

		defer close(stream)
		if true {
			return
		}
		stream <- "hi specter"
	}()

	/*
		Notice how the loop doesn't need an exit criteria and the range doesn't return the second
		boolean value. The specifics of handling a closed channel are managed for you to keep the
		loop concise


		Closing a channel is also one of the ways you can signal multiple goroutines simultaneously.
		if you have n goroutine waiting on a single channel, instead of writing n times to the channel
		to unblock each goroutine, you can simply close the goroutine
		check v6 for this.
	*/
	for val := range stream {
		fmt.Println("received ", val)
	}

}
