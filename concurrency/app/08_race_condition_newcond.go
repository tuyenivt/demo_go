package app

import (
	"fmt"
	"sync"
)

var raceConditionNewCondShareValue = 3000
var raceConditionNewCondWaitGroup = sync.WaitGroup{}
var raceConditionNewCondMutex = sync.Mutex{}
var raceConditionNewCondNewCond = sync.NewCond(&raceConditionNewCondMutex)

func RaceConditionNewCond() {
	fmt.Println("raceConditionNewCondShareValue start value = ", raceConditionNewCondShareValue)
	raceConditionNewCondWaitGroup.Add(2)
	go raceConditionNewCondAdd()
	go raceConditionNewCondSubtract()
	raceConditionNewCondWaitGroup.Wait()
	fmt.Println("raceConditionNewCondShareValue end value = ", raceConditionNewCondShareValue)
}

func raceConditionNewCondAdd() {
	for i := 0; i < 1000; i++ {
		raceConditionNewCondMutex.Lock()
		raceConditionNewCondShareValue += 100
		raceConditionNewCondNewCond.Signal()
		raceConditionNewCondMutex.Unlock()
	}
	raceConditionNewCondWaitGroup.Done()
}

func raceConditionNewCondSubtract() {
	for i := 0; i < 1000; i++ {
		raceConditionNewCondMutex.Lock()
		for raceConditionNewCondShareValue-100 < 0 {
			raceConditionNewCondNewCond.Wait()
		}
		raceConditionNewCondShareValue -= 100
		raceConditionNewCondMutex.Unlock()
	}
	raceConditionNewCondWaitGroup.Done()
}
