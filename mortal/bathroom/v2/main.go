// local join point for EnterBathroom

package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	PartyDemocrat   = "democrat"
	PartyRepublican = "republican"
)

type Bathroom struct {
	party         string
	maxCapacity   int
	currentLength int
	sync.Mutex
	// cond is a rendezvous point for goroutines waiting for or announcing the occurence of an event
	// an event is an arbitrary signal between two or more goroutines that carry no information other
	// other than the face that it has occurred
	signal *sync.Cond
}

func NewBathroom(capacity int) *Bathroom {
	b := &Bathroom{
		party:         "",
		maxCapacity:   capacity,
		currentLength: 0,
	}
	b.signal = sync.NewCond(&b.Mutex)

	return b
}

func EnterBathroom(wg *sync.WaitGroup, b *Bathroom, name string, party string) {
	use := func(wg *sync.WaitGroup, b *Bathroom, name string, party string) {
		// added wg(join point) in order to return full response, otherwise main exits once all
		// have entered bathroom and we won't see in stdout `finished using bathroom` for some folks
		defer wg.Done()
		time.Sleep(time.Second)
		b.Lock()
		b.currentLength--
		fmt.Printf("%s of type %s finished using bathroom\n", name, party)
		b.Unlock()
		b.signal.Broadcast()
	}

	var localWg sync.WaitGroup

	b.Lock()
	// signal on the condition doesn't necessarily mean what you have been waiting for has
	// occurred-only that something has occurred
	for b.currentLength > 0 && b.party != party {
		fmt.Printf("%s of party %s is waiting to use the bathroom\n", name, party)
		b.signal.Wait()
		// this is a blocking call and the goroutine will be suspended too, allowing other
		// goroutines to run on os thread
		// on entering wait , unlock is called on cond variable's locker
		// on exiting wait, lock is called on cond variable's locker
	}

	if b.currentLength < b.maxCapacity {
		fmt.Printf("%s of party %s entered the bathroom\n", name, party)
		b.currentLength++
		b.party = party
		b.Unlock()
		localWg.Add(1)
		go use(&localWg, b, name, party)
		localWg.Wait()
		return
	}

	for b.currentLength < b.maxCapacity {
		b.signal.Wait()
	}
	b.currentLength++
	b.party = party
	localWg.Add(1)
	go use(&localWg, b, name, party)
	b.Unlock()
	localWg.Wait()
}

func DemocratUseBathroom(wg *sync.WaitGroup, b *Bathroom, name string) {
	defer wg.Done()
	EnterBathroom(wg, b, name, PartyDemocrat)
}

func RepublicanUseBathroom(wg *sync.WaitGroup, b *Bathroom, name string) {
	defer wg.Done()
	EnterBathroom(wg, b, name, PartyRepublican)
}

func main() {
	var wg sync.WaitGroup

	b := NewBathroom(3)
	wg.Add(4)
	// waitgroup is a great way to wait for a set of concurrent operations to complete when you
	// either don't care about the results of the concurrent operation, or you have other means of
	// collecting their results. If neither of these conditions are true, use channels and a
	// select statement instead
	go RepublicanUseBathroom(&wg, b, "RA")
	go DemocratUseBathroom(&wg, b, "DA")
	go DemocratUseBathroom(&wg, b, "DB")
	go DemocratUseBathroom(&wg, b, "DC")

	// a race condition occurs when two or more operations must execute in the correct order,
	// but the program has not been written so that this order is guaranteed to be maintained
	wg.Wait() // creates a join point
	// Join points are what guarantee our program's correctness and remove the race condition
}
