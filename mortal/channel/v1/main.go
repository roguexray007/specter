package main

import "fmt"

func main() {

	go func() {
		fmt.Println("hi specter")
	}()

	fmt.Println("main exiting")

	// very important
	// just because a goroutine is scheduled, there was no guarantee that it would run before the main
	// exited
}
