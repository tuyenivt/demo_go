package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

func main() {
	go hello("World")
	time.Sleep(time.Second)

	simpleWaitGroup()
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
