package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	doWork := func(wg *sync.WaitGroup, start *sync.Cond, id int) {
		// very important. try removing localwg and it will fail or try commenting localwg wait
		var localwg sync.WaitGroup
		localwg.Add(1)
		go func() {
			defer func() {
				wg.Done()
				fmt.Println("goroutine ", id, " exited")
			}()
			localwg.Done()
			i := 0
			for i != 2 {
				start.L.Lock()
				start.Wait()
				fmt.Printf("gorutine no %d completed it's work in broadcast term %d\n", id, i)
				i++

				start.L.Unlock()

			}
		}()
		/*

			The localwg is a sync.WaitGroup created within each (doWork function) to ensure that the
			main work loop doesn't begin until the goroutine is fully initialized.
			Each goroutine calls localwg.Add(1) to increment its internal counter and localwg.Done() to
			decrement it immediately after starting the goroutine.
			After the main setup (increment and immediate decrement of localwg), the goroutine
			proceeds to wait on the shared start condition variable.
			Only when all goroutines have executed localwg.Done() (which occurs right after they've initialized
			themselves) does the main loop reach localwg.Wait(), allowing the main loop to proceed and call
			start.Broadcast().
			So, indeed, localwg is ensuring that all goroutines have properly initialized themselves before
			the main loop proceeds to signal the start condition variable with start.Broadcast().
		*/
		localwg.Wait()
	}

	var start *sync.Cond
	start = sync.NewCond(&sync.Mutex{})

	var wg sync.WaitGroup
	broadcastCounter := 2

	for i := 0; i < 3; i++ {
		wg.Add(1)
		doWork(&wg, start, i)
	}
	fmt.Println("unblocking goroutines")
	for i := 0; i < broadcastCounter; i++ {
		start.Broadcast()
		time.Sleep(time.Second) // allow some time for work to be completed before sending another event
	}

	wg.Wait() // created a join point
}
