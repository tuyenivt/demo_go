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

	simpleChannel()
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
