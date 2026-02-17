package main

import (
	"fmt"
	"sync"
)

// 1. PRODUCER: Owns the 'tasks' channel
func produce(done <-chan struct{}, numbers ...int) <-chan int {
	tasks := make(chan int)
	go func() {
		defer close(tasks)
		for _, n := range numbers {
			select {
			case <-done:
				return
			case tasks <- n:
			}
		}
	}()
	return tasks
}

// 2. WORKER: Transforms data (Fan-out stage)
func worker(done <-chan struct{}, tasks <-chan int) <-chan int {
	results := make(chan int)
	go func() {
		defer close(results)
		for n := range tasks {
			select {
			case <-done:
				return
			case results <- n * n: // Squaring the number as "work"
			}
		}
	}()
	return results
}

// 3. MERGER: THE FAN-IN OWNER
// This function takes multiple result channels and merges them into one.
func fanIn(done <-chan struct{}, channels ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	multiplexedStream := make(chan int)

	// Helper function to drain a single channel into the multiplexed stream
	multiplex := func(c <-chan int) {
		defer wg.Done()
		for n := range c {
			select {
			case <-done:
				return
			case multiplexedStream <- n:
			}
		}
	}

	wg.Add(len(channels))
	for _, c := range channels {
		go multiplex(c)
	}

	// Wait for all workers to finish, then close the stream we OWN
	go func() {
		wg.Wait()
		close(multiplexedStream)
	}()

	return multiplexedStream
}

// 4. COORDINATOR: Orchestrates the entire pipeline
func main() {
	done := make(chan struct{})
	defer close(done)

	// Stage 1: Produce tasks
	tasks := produce(done, 1, 2, 3, 4, 5, 6)

	// Stage 2: Fan-out (3 workers)
	w1 := worker(done, tasks)
	w2 := worker(done, tasks)
	w3 := worker(done, tasks)

	// Stage 3: Fan-in (Merge 3 worker channels into 1)
	finalResults := fanIn(done, w1, w2, w3)

	// Stage 4: Consume final results
	fmt.Println("Coordinator: Receiving merged results...")
	for res := range finalResults {
		fmt.Printf("Result: %d\n", res)
	}
	fmt.Println("Coordinator: Pipeline complete.")
}
