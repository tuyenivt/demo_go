package app

import (
	"fmt"
	"strconv"
)

func SimpleChannel() {
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
