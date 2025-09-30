package app

import (
	"fmt"
	"strconv"
	"sync"
)

func SimpleWaitGroup() {
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
