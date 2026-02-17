package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	// 1. The Locker and the Condition
	// The Mutex protects the shared state (the slice)
	mutex := &sync.Mutex{}
	cond := sync.NewCond(mutex)

	queue := make([]interface{}, 0)

	// 2. THE CONSUMER PATTERN
	removeFromQueue := func(id int, wg *sync.WaitGroup) {
		defer wg.Done()

		cond.L.Lock()
		// Rule: Always use a for-loop to check the condition
		for len(queue) == 0 {
			fmt.Printf("Consumer %d: Queue empty, going to sleep...\n", id)
			cond.Wait() // Atomically unlocks and suspends
		}

		// When Wait() returns, the lock has been re-acquired
		item := queue[0]
		queue = queue[1:]
		fmt.Printf("Consumer %d: Removed %v. Remaining: %d\n", id, item, len(queue))

		cond.L.Unlock()
	}

	// 3. STARTING CONSUMERS
	var wg sync.WaitGroup
	for i := 1; i <= 2; i++ {
		wg.Add(1)
		go removeFromQueue(i, &wg)
	}

	// 4. THE PRODUCER PATTERN
	time.Sleep(1 * time.Second) // Wait to show consumers sleeping

	cond.L.Lock()
	fmt.Println("Producer: Adding items to queue...")
	queue = append(queue, "Task A", "Task B")

	// Rule: Use Signal() to wake one, or Broadcast() to wake all
	fmt.Println("Producer: Signaling all consumers (Broadcast)...")
	cond.Broadcast()
	cond.L.Unlock()

	wg.Wait()
	fmt.Println("Main: All tasks processed.")
}
