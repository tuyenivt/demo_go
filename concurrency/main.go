package main

import (
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	go hello("World")
	time.Sleep(time.Second)

	simpleWaitGroup()

	simpleChannel()

	multiChannel()

	raceConditionMutex()

	raceConditionAtomic()

	raceConditionNewCond()
}

func hello(name string) {
	fmt.Printf("Hello, %s\n", name)
}

func simpleWaitGroup() {
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go wgHello(strconv.Itoa(i), &wg)
	}
	wg.Wait()
}

func wgHello(name string, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("WaitGroup Hello, %s\n", name)
}

func simpleChannel() {
	ch := make(chan string, 5)
	defer close(ch)
	for i := 0; i < 5; i++ {
		go chHello(strconv.Itoa(i), ch)
	}
	for i := 0; i < 5; i++ {
		fmt.Print(<-ch)
	}
}

func chHello(name string, ch chan string) {
	ch <- fmt.Sprintf("Channel Hello, %s\n", name)
}

func multiChannel() {
	ch1 := make(chan string)
	ch2 := make(chan string)
	ch3 := make(chan string)
	defer close(ch1)
	defer close(ch2)
	defer close(ch3)

	go multiHello(1, ch1)
	go multiHello(2, ch2)
	go multiHello(3, ch3)

Loop:
	for {
		select {
		case msg := <-ch1:
			fmt.Print("MultiChannel 1 - " + msg)
		case msg := <-ch2:
			fmt.Print("MultiChannel 2 - " + msg)
		case msg := <-ch3:
			fmt.Print("MultiChannel 3 - " + msg)
			break Loop
		}
	}
}

func multiHello(name int, ch chan string) {
	time.Sleep(time.Duration(name) * time.Millisecond)
	ch <- fmt.Sprintf("MultiChannel Hello, %d\n", name)
}

var raceConditionMutexShareValue = 1000
var raceConditionMutexWaitGroup = sync.WaitGroup{}
var raceConditionMutexMutex = sync.Mutex{}

func raceConditionMutex() {
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

var raceConditionAtomicShareValue int32 = 2000
var raceConditionAtomicWaitGroup = sync.WaitGroup{}

func raceConditionAtomic() {
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

var raceConditionNewCondShareValue = 3000
var raceConditionNewCondWaitGroup = sync.WaitGroup{}
var raceConditionNewCondMutex = sync.Mutex{}
var raceConditionNewCondNewCond = sync.NewCond(&raceConditionNewCondMutex)

func raceConditionNewCond() {
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
