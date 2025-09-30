package app

import (
	"context"
	"fmt"
	"time"
)

func SimpleContextTimeout() {
	fmt.Println("simpleContextTimeout start")
	time.Sleep(1 * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	go simpleContextTimeoutHello(ctx, "World")

	time.Sleep(1 * time.Second)
	fmt.Println("simpleContextTimeout end")
}

func simpleContextTimeoutHello(ctx context.Context, name string) {
	fmt.Println("simpleContextTimeout processing...")
	select {
	case <-time.After(2 * time.Second):
		fmt.Printf("simpleContextTimeout processed %s\n", name)
	case <-ctx.Done():
		fmt.Println("simpleContextTimeout cancel result:", ctx.Err())
	}
}
