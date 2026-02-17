package main

import (
	"fmt"
	"time"
)

// --- THE BRIDGE ---
// It owns 'valStream'. It consumes a stream of channels and
// flattens them into that single 'valStream'.
func bridge(done <-chan struct{}, chanStream <-chan (<-chan int)) <-chan int {
	valStream := make(chan int)
	go func() {
		defer close(valStream)
		for {
			var stream <-chan int
			select {
			case maybeStream, ok := <-chanStream:
				if !ok {
					return
				}
				stream = maybeStream
			case <-done:
				return
			}

			// Draining the inner channel. We use the Or-Done logic
			// here to ensure we don't hang if 'done' is closed.
			for val := range stream {
				select {
				case <-done:
					return
				case valStream <- val:
				}
			}
		}
	}()
	return valStream
}

// --- HELPER: GEN-STREAM ---
// This simulates the "Master Producer" that creates other channels.
func genStream(done <-chan struct{}) <-chan (<-chan int) {
	chanStream := make(chan (<-chan int))
	go func() {
		defer close(chanStream)
		for i := 0; i < 3; i++ {
			// We create a new sub-channel for every iteration
			subChan := make(chan int)
			go func(v int) {
				defer close(subChan)
				for j := 0; j < 3; j++ {
					subChan <- v*10 + j // Sends 0,1,2 then 10,11,12 etc.
					time.Sleep(100 * time.Millisecond)
				}
			}(i)

			select {
			case <-done:
				return
			case chanStream <- subChan:
			}
		}
	}()
	return chanStream
}

// --- THE COORDINATOR ---
func main() {
	done := make(chan struct{})
	defer close(done)

	fmt.Println("Starting Bridge Pipeline...")

	// 1. Get our stream of channels
	nestedChannels := genStream(done)

	// 2. Bridge them into a single channel
	singleStream := bridge(done, nestedChannels)

	// 3. Consume the results simply
	// Notice: The consumer doesn't care that the data came from
	// 3 different internal channels.
	for val := range singleStream {
		fmt.Printf("Received: %d\n", val)
	}

	fmt.Println("Bridge Pipeline Complete.")
}
