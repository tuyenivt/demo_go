package app

import (
	"fmt"
	"sync"
	"time"
)

type Job struct {
	ID   int
	Desc string
}

type Result struct {
	JobID  int
	Output string
}

func SimpleWorkerPool() {
	jobs := make(chan Job, 10)
	results := make(chan Result, 10)

	// worker pool
	worker_count := 3
	var wg sync.WaitGroup
	wg.Add(worker_count)
	for i := 1; i <= worker_count; i++ {
		go doJob(i, &wg, jobs, results)
	}

	// cleanup
	go func() {
		wg.Wait()
		close(results)
		fmt.Println("All jobs finished. Signal close results channel sent.")
	}()

	// send job to jobs channel
	job_count := 100
	go createJob(job_count, jobs)

	// receive output from results channel
	getResult(job_count, results)
}

func createJob(job_count int, jobs chan<- Job) {
	for i := 1; i <= job_count; i++ {
		time.Sleep(time.Millisecond * 100)
		jobs <- Job{ID: i, Desc: fmt.Sprintf("Job %d", i)}
	}
	close(jobs)
	fmt.Println("All jobs created. Signal close jobs channel sent.")
}

func doJob(workerID int, wg *sync.WaitGroup, jobs <-chan Job, results chan<- Result) {
	defer wg.Done()
	for job := range jobs {
		time.Sleep(time.Millisecond * 500)
		results <- Result{JobID: job.ID, Output: fmt.Sprintf("Worker %d finished job id %d with desc %s", workerID, job.ID, job.Desc)}
	}
}

func getResult(job_count int, results <-chan Result) {
	for i := 1; i <= job_count; i++ {
		fmt.Println(<-results)
	}
}
