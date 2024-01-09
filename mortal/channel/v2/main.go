package main

import "fmt"

func main() {
	var stream chan string
	stream = make(chan string)

	go func() {
		stream <- "hi specter"
	}()

	val := <-stream
	fmt.Println("recieved ", val)

	/*
		we highlighted the fact that just because a goroutine was scheduled, there was no guarantee
		that it would run before the process exited;yet this example is complete and correct with no
		code omitted. You may have been wondering why the anonymous goroutine completes before
		the main goroutine does;or did i just get lucky when i ran this?

		This works because channels in go are blocking. This means that any goroutine that attempts to
		write to a channel that is full will wait until the channel has been emptied, and
		any goroutine that attempts to read from a channel that is empty will wait until at least
		one item is placed on it.
	*/
}
