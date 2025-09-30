package app

import (
	"fmt"
	"time"
)

func SimpleGoroutine(name string) {
	go func(name string) {
		fmt.Printf("Hello, %s\n", name)
	}(name)
	time.Sleep(100 * time.Millisecond)
}
