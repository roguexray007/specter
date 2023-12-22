package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup

func sayHello() {
	defer wg.Done()
	fmt.Println("hello from function ")
}

func main() {

	wg.Add(1)
	// call via function Name
	go sayHello()

	// call via anonymous function
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("hello from anonymous function")
	}() //we must invoke the anonymous function immediately to use the go keyword

	// alternatively we can assign the function to a variable and call the anonymous func
	// like this
	wg.Add(1)
	yoSayHello := func() {
		defer wg.Done()
		fmt.Println("hello from function assigned to a varaiable")
	}

	go yoSayHello()

	wg.Wait()
}
