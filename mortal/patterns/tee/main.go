package main

import (
	"fmt"
)

// --- THE TEE PATTERN ---
func tee(done <-chan struct{}, in <-chan int) (<-chan int, <-chan int) {
	out1 := make(chan int)
	out2 := make(chan int)

	go func() {
		defer close(out1)
		defer close(out2)

		// We use our friend orDone to ensure we can exit while reading 'in'
		for val := range in {
			// We create local copies of the channels so we can
			// use 'nil' to track which ones have received the current value
			var out1, out2 = out1, out2

			// This inner loop ensures both channels get the value
			// before we pull the next value from 'in'
			for i := 0; i < 2; i++ {
				select {
				case <-done:
					return
				case out1 <- val:
					out1 = nil // Setting to nil disables this select case
				case out2 <- val:
					out2 = nil // Setting to nil disables this select case
				}
			}
		}
	}()
	return out1, out2
}

func main() {
	done := make(chan struct{})
	defer close(done)

	// 1. Source of data
	source := make(chan int)
	go func() {
		for i := 1; i <= 3; i++ {
			source <- i
		}
		close(source)
	}()

	// 2. Split the stream
	logChan, saveChan := tee(done, source)

	// 3. Consume from both (using a WaitGroup to wait for both)
	go func() {
		for v := range logChan {
			fmt.Printf("[LOG]: %d\n", v)
		}
	}()

	for v := range saveChan {
		fmt.Printf("[SAVE]: %d to database\n", v)
	}
}
