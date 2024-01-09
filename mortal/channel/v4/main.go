package main

import "fmt"

func main() {
	var stream chan string
	stream = make(chan string)

	go func() {
		// cause deadlock, as main is waiting for value, but doesn;t get one, this gorutine exits
		// question is how do we prevent deadlock like this?
		// ans- close channel
		defer close(stream)
		if true {
			return
		}
		stream <- "hi specter"
	}()

	/*
		The receiving form of the <- operator can also optionally return 2 values
		What does the boolean signify? The second return value is a way for a read opeartion
		to indicate whether the read off the channel was a value generated by a write elsewhere
		or a default value generated from a closed channel.

		Wait a second , what's a closed channel??
		In programs, it's very useful to be able to indicate that no more values will be sent
		over a channel. This helps downstream processes know when to move on, exit, reopen comm
		on a new or different channel etc. So closing a channel is like a special sentinel
		that says "hey, upstream isn't going to be writing any more values, do what you will"

		we can read on a closed channel, and infact we can continue reading on this channel
		indefinitely despite the channel remaining closed. This is to allow support for multiple
		downstreams read from a single upstream writer on channel.

		Second return value opens up a new pattern using range over a channel. check v5
	*/
	val, ok := <-stream
	fmt.Println("recieved ", val, ok)

}