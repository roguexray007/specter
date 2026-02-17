package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

// 1. PRODUCER PATTERN: Returns a channel it owns and an error-handling closure
func startService(ctx context.Context, name string, delay time.Duration) (<-chan string, func() error) {
	out := make(chan string)

	// The closure that will be executed by the errgroup
	run := func() error {
		defer close(out) // Owner closes the channel
		for i := 1; i <= 3; i++ {
			// Simulate a failure for Service B
			if name == "Service B" && i == 2 {
				return errors.New("Service B encountered a critical error")
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			case out <- fmt.Sprintf("%s: Data %d", name, i):
				time.Sleep(delay)
			}
		}
		return nil
	}

	return out, run
}

// 2. FAN-IN PATTERN: Merges multiple streams into one owned stream
func merge(ctx context.Context, channels ...<-chan string) <-chan string {
	merged := make(chan string)
	var wg sync.WaitGroup // Using WG here to track the "draining" goroutines

	multiplex := func(c <-chan string) {
		defer wg.Done()
		for v := range c {
			select {
			case <-ctx.Done():
				return
			case merged <- v:
			}
		}
	}

	wg.Add(len(channels))
	for _, c := range channels {
		go multiplex(c)
	}

	// Ownership Duty: Background closer
	go func() {
		wg.Wait()
		close(merged)
	}()

	return merged
}

// 3. COORDINATOR PATTERN: Uses ErrGroup for lifecycle management
func main() {
	g, ctx := errgroup.WithContext(context.Background())

	// Initialize Services
	chA, runA := startService(ctx, "Service A", 100*time.Millisecond)
	chB, runB := startService(ctx, "Service B", 150*time.Millisecond)
	chC, runC := startService(ctx, "Service C", 200*time.Millisecond)

	// Register runners with the group
	g.Go(runA)
	g.Go(runB)
	g.Go(runC)

	// Use Fan-in to consolidate results
	finalResults := merge(ctx, chA, chB, chC)

	// Consume results in a separate goroutine to avoid blocking g.Wait()
	go func() {
		for res := range finalResults {
			fmt.Println("Result:", res)
		}
	}()

	// Wait for completion or error
	if err := g.Wait(); err != nil {
		fmt.Printf("\n[SYSTEM] Execution stopped: %v\n", err)
	} else {
		fmt.Println("\n[SYSTEM] All services completed successfully.")
	}
}

/*
	Key Takeaways from this Pattern
The Run Closure: Notice how startService returns both the channel and a function. This allows the producer to define its logic, but gives the errgroup in main the power to trigger it.

Automatic Short-Circuit: Because we used errgroup.WithContext, the moment Service B returns an error, the ctx passed to Service A and Service C is cancelled. They will hit their case <-ctx.Done() and exit immediately.

Clean Drains: Even if an error occurs, the merge function’s WaitGroup ensures we don't close the finalResults channel until all sub-goroutines have stopped trying to write to it.
*/
