package main

import (
	"fmt"
	"sync"
	"time"
)

// --- THE PRODUCER ---
// Owns the data stream.
func produceTasks(count int) <-chan int {
	tasks := make(chan int)
	go func() {
		defer close(tasks) // Owner's duty: close when done
		for i := 1; i <= count; i++ {
			tasks <- i
		}
	}()
	return tasks
}

// --- THE WORKER ---
// Does not own the channel; just consumes until it's closed.
func worker(id int, tasks <-chan int, wg *sync.WaitGroup) {
	defer wg.Done() // Signal that this specific worker is finished
	for t := range tasks {
		// Simulate a time-consuming task
		fmt.Printf("Worker %d: Processing task %d\n", id, t)
		time.Sleep(200 * time.Millisecond)
	}
	fmt.Printf("Worker %d: Finished and exiting\n", id)
}

func main() {
	// 1. Start the producer (10 tasks)
	taskStream := produceTasks(10)

	// 2. The WaitGroup tracks our workers
	var wg sync.WaitGroup

	// 3. FAN-OUT: Start 3 worker goroutines
	numWorkers := 3
	fmt.Printf("Starting %d workers to process tasks...\n", numWorkers)

	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, taskStream, &wg)
	}

	// 4. Wait for all workers to finish their loops
	wg.Wait()
	fmt.Println("Main: All tasks complete. Pipeline shut down.")
}
