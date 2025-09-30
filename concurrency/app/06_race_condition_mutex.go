package app

import (
	"fmt"
	"sync"
)

var raceConditionMutexShareValue = 1000
var raceConditionMutexWaitGroup = sync.WaitGroup{}
var raceConditionMutexMutex = sync.Mutex{}

func RaceConditionMutex() {
	fmt.Println("raceConditionMutexShareValue start value = ", raceConditionMutexShareValue)
	raceConditionMutexWaitGroup.Add(2)
	go raceConditionMutexAdd()
	go raceConditionMutexSubtract()
	raceConditionMutexWaitGroup.Wait()
	fmt.Println("raceConditionMutexShareValue end value = ", raceConditionMutexShareValue)
}

func raceConditionMutexAdd() {
	for i := 0; i < 1000; i++ {
		raceConditionMutexMutex.Lock()
		raceConditionMutexShareValue += 100
		raceConditionMutexMutex.Unlock()
	}
	raceConditionMutexWaitGroup.Done()
}

func raceConditionMutexSubtract() {
	for i := 0; i < 1000; i++ {
		raceConditionMutexMutex.Lock()
		raceConditionMutexShareValue -= 100
		raceConditionMutexMutex.Unlock()
	}
	raceConditionMutexWaitGroup.Done()
}
