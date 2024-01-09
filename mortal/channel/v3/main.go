package main

import "fmt"

func main() {
	var stream chan string
	stream = make(chan string)

	go func() {
		// cause deadlock, as main is waiting for value, but doesn;t get one, this gorutine exits
		// question is how do we prevent deadlock like this?
		if true {
			return
		}
		stream <- "hi specter"
	}()

	val := <-stream
	fmt.Println("recieved ", val)

}
