package app

import (
	"fmt"
	"sync"
	"sync/atomic"
)

var raceConditionAtomicShareValue int32 = 2000
var raceConditionAtomicWaitGroup = sync.WaitGroup{}

func RaceConditionAtomic() {
	fmt.Println("raceConditionAtomicShareValue start value = ", raceConditionAtomicShareValue)
	raceConditionAtomicWaitGroup.Add(2)
	go raceConditionAtomicAdd()
	go raceConditionAtomicSubtract()
	raceConditionAtomicWaitGroup.Wait()
	fmt.Println("raceConditionAtomicShareValue end value = ", raceConditionAtomicShareValue)
}

func raceConditionAtomicAdd() {
	for i := 0; i < 1000; i++ {
		atomic.AddInt32(&raceConditionAtomicShareValue, 100)
	}
	raceConditionAtomicWaitGroup.Done()
}

func raceConditionAtomicSubtract() {
	for i := 0; i < 1000; i++ {
		atomic.AddInt32(&raceConditionAtomicShareValue, -100)
	}
	raceConditionAtomicWaitGroup.Done()
}
