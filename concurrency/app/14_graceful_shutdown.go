package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func SimpleGracefulShutdown() {
	// Create root context that cancels on SIGINT or SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	// WaitGroup for running goroutines
	var wg sync.WaitGroup

	// Open a file
	file, err := os.Create("temp.log")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		fmt.Println("Closing file...")
		file.Close()
	}()

	// Start a background task (simulated scheduler)
	wg.Add(1)
	go func() {
		defer wg.Done()
		runScheduler(ctx)
	}()

	// Simulate another worker
	wg.Add(1)
	go func() {
		defer wg.Done()
		doWork(ctx)
	}()

	// Wait for shutdown signal
	<-ctx.Done()
	fmt.Println("\nShutdown signal received...")

	// Optional: give workers some time to finish gracefully
	graceCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		fmt.Println("All tasks completed.")
	case <-graceCtx.Done():
		fmt.Println("Timeout reached. Forcing shutdown.")
	}
}

func runScheduler(ctx context.Context) {
	fmt.Println("Scheduler started.")
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Scheduler stopping...")
			return
		case t := <-ticker.C:
			fmt.Println("Scheduled task at", t)
		}
	}
}

func doWork(ctx context.Context) {
	fmt.Println("Worker started.")
	for i := 0; i < 5; i++ {
		select {
		case <-ctx.Done():
			fmt.Println("Worker stopping early...")
			return
		default:
			fmt.Printf("Working... %d\n", i)
			time.Sleep(1 * time.Second)
		}
	}
	fmt.Println("Worker done.")
}
