package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// --- PRODUCER (Owner of 'tasks') ---
func produce(ctx context.Context) <-chan int {
	tasks := make(chan int)
	go func() {
		defer close(tasks)
		for i := 1; ; i++ {
			select {
			case <-ctx.Done(): // Owner signals stop
				return
			case tasks <- i:
				time.Sleep(200 * time.Millisecond)
			}
		}
	}()
	return tasks
}

// --- WORKER (Background Goroutine) ---
func worker(ctx context.Context, id int, tasks <-chan int, wg *sync.WaitGroup) {
	defer wg.Done() // Signal to coordinator that this worker is finished
	fmt.Printf("Worker %d: Started\n", id)

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Worker %d: Context cancelled, exiting...\n", id)
			return
		case val, ok := <-tasks:
			if !ok {
				return
			}
			// Simulate processing
			fmt.Printf("Worker %d: Processed %d\n", id, val)
		}
	}
}

// --- COORDINATOR (The Authority) ---
func main() {
	// 1. Setup Ownership/Cancellation
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// 2. Setup Synchronization
	var wg sync.WaitGroup

	// 3. Start Producer
	taskStream := produce(ctx)

	// 4. FAN-OUT: Start multiple worker goroutines
	numWorkers := 3
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(ctx, i, taskStream, &wg)
	}

	// 5. Wait for the workers to finish their cleanup
	wg.Wait()
	fmt.Println("Main: All goroutines finished cleanly.")
}
