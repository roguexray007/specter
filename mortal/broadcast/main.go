package main

import (
	"fmt"
	"sync"
)

type Button struct {
	clickSignal *sync.Cond
}

func NewButton() *Button {
	b := &Button{}
	b.clickSignal = sync.NewCond(&sync.Mutex{})

	return b
}

type Handler func()

func main() {
	RegisterHandler := func(b *Button, fn Handler) {
		var localWg sync.WaitGroup
		localWg.Add(1)
		go func() {
			localWg.Done() // don't defer  -  unblock registerHandler
			b.clickSignal.L.Lock()
			defer b.clickSignal.L.Unlock()
			b.clickSignal.Wait()
			fn()
		}()
		localWg.Wait()
	}

	b := NewButton()

	var wg sync.WaitGroup

	wg.Add(3)

	RegisterHandler(b, func() {
		defer wg.Done()
		fmt.Println("handler 1 running")
	})

	RegisterHandler(b, func() {
		defer wg.Done()
		fmt.Println("handler 2 running")
	})

	RegisterHandler(b, func() {
		defer wg.Done()
		fmt.Println("handler 3 running")
	})

	b.clickSignal.Broadcast()

	wg.Wait()
}
