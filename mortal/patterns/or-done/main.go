package main

import (
	"fmt"
	"time"
)

// --- THE OR-DONE WRAPPER ---
// This function takes a channel you DON'T own and returns a channel you DO own.
func orDone(done <-chan struct{}, dataStream <-chan int) <-chan int {
	valStream := make(chan int)
	go func() {
		defer close(valStream) // We own this, so we close it
		for {
			select {
			case <-done:
				return // Exit if we get the done signal
			case v, ok := <-dataStream:
				if !ok {
					return // Exit if the original channel is closed
				}
				// Pass the value to the consumer, but stay alert for 'done'
				select {
				case <-done:
					return
				case valStream <- v:
				}
			}
		}
	}()
	return valStream
}

func main() {
	done := make(chan struct{})

	// Simulate a producer we don't control
	slowProducer := make(chan int)
	go func() {
		defer close(slowProducer)
		for i := 1; i <= 5; i++ {
			time.Sleep(1 * time.Second) // Very slow!
			slowProducer <- i
		}
	}()

	// We use the Or-Done pattern to consume safely
	fmt.Println("Starting to consume...")

	// We only want to wait for 2 seconds total
	go func() {
		time.Sleep(2 * time.Second)
		fmt.Println("\n[Coordinator] Time's up! Closing done channel...")
		close(done)
	}()

	// Use the wrapped channel in the loop
	for val := range orDone(done, slowProducer) {
		fmt.Printf("Received: %d ", val)
	}

	fmt.Println("\nMain: Cleanly exited without leaking.")
}
