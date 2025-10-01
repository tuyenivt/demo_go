package app

import (
	"fmt"
	"sync"
)

func SimpleSyncOnce() {
	var wg sync.WaitGroup

	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go simpleSyncOnceTask(i, &wg)
	}

	wg.Wait()
}

var once sync.Once // sync.Once used for singleton, load config, create DB connection, ...

func simpleSyncOnceTask(i int, wg *sync.WaitGroup) {
	defer wg.Done()

	once.Do(simpleSyncOnceInitial)

	fmt.Printf("Processing task %d\n", i)
}

func simpleSyncOnceInitial() {
	fmt.Println("Expected run once...")
}
