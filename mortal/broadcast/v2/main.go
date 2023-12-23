package main

import (
	"fmt"
	"sync"
	"time"
)

type Button struct {
	clickSignal *sync.Cond
	done        bool
}

func NewButton() *Button {
	b := &Button{}
	b.clickSignal = sync.NewCond(&sync.Mutex{})

	return b
}

type Handler func()

func main() {
	var wg sync.WaitGroup

	RegisterHandler := func(wg *sync.WaitGroup, b *Button, fn Handler) {
		var localWg sync.WaitGroup
		localWg.Add(1)
		go func() {
			defer wg.Done()
			localWg.Done() // don't defer  -  unblock registerHandler
			b.clickSignal.L.Lock()
			for !b.done {
				b.clickSignal.Wait()
				if !b.done {
					fn()
				}
			}
			b.clickSignal.L.Unlock()
		}()
		localWg.Wait()
	}

	b := NewButton()

	wg.Add(3) // no of registerHandler

	RegisterHandler(&wg, b, func() {
		fmt.Println("handler 1 running")
	})

	RegisterHandler(&wg, b, func() {
		fmt.Println("handler 2 running")
	})

	RegisterHandler(&wg, b, func() {
		fmt.Println("handler 3 running")
	})

	for i := 0; i < 3; i++ {
		b.clickSignal.Broadcast()
		time.Sleep(time.Second)
	}
	b.done = true
	b.clickSignal.Broadcast()
	wg.Wait()
}
