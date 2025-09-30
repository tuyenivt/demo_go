package app

import (
	"context"
	"fmt"
	"time"
)

func SimpleContextCancel() {
	fmt.Println("simpleContextCancel start")
	time.Sleep(1 * time.Second)

	ctx, cancel := context.WithCancel(context.Background())

	go simpleContextCancelHello(ctx, "World")

	time.Sleep(500 * time.Millisecond)
	fmt.Println("simpleContextCancel cancelling context...")
	cancel()

	time.Sleep(1 * time.Second)
	fmt.Println("simpleContextCancel end")
}

func simpleContextCancelHello(ctx context.Context, name string) {
	i := 0
	for {
		i++
		select {
		case <-ctx.Done():
			fmt.Println("simpleContextCancel cancel result:", ctx.Err())
			return
		default:
			fmt.Printf("simpleContextCancel %s processing %d...\n", name, i)
			time.Sleep(100 * time.Millisecond)
		}
	}
}
