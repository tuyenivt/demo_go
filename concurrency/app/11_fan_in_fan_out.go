package app

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func SimpleFanInFanOut() {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	dataStream := simpleFanInFanOutDataGenerator(10, rnd)

	worker1 := simpleFanInFanOutWorker(1, dataStream)
	worker2 := simpleFanInFanOutWorker(2, dataStream)

	merged := simpleFanInFanOutFanIn(worker1, worker2)

	for result := range merged {
		fmt.Println("Result:", result)
	}
}

func simpleFanInFanOutDataGenerator(count int, rnd *rand.Rand) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for i := 0; i < count; i++ {
			num := rnd.Intn(100)
			out <- num
		}
	}()
	return out
}

func simpleFanInFanOutWorker(id int, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			fmt.Printf("Worker %d processing %d\n", id, n)
			time.Sleep(time.Millisecond * 100)
			out <- n * 2
		}
	}()
	return out
}

func simpleFanInFanOutFanIn(channels ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	output := func(c <-chan int) {
		defer wg.Done()
		for n := range c {
			out <- n
		}
	}

	wg.Add(len(channels))
	for _, ch := range channels {
		go output(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
