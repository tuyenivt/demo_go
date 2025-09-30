package app

import (
	"fmt"
	"time"
)

func MultiChannel() {
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
		default:
			fmt.Println("No data from any channel...")
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func multiHello(name int, ch chan string) {
	time.Sleep(time.Duration(name) * time.Millisecond)
	ch <- fmt.Sprintf("MultiChannel Hello, %d\n", name)
}
