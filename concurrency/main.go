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

	multiChannel()
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
