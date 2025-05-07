package extractor

import (
	"fmt"
	"sync"
)

// Job contains data about Page processed and current page number
type Job struct {
	Page    Page
	PageNum int
}

// JobErr stores error for workerID if occurs
type JobErr struct {
	Err      error
	WorkerID int
}

// Worker process page extraction
func Worker(id int, jobs <-chan Job, errors chan<- JobErr, done chan<- struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	// panic recovery
	defer func() {
		if r := recover(); r != nil {
			errors <- JobErr{Err: fmt.Errorf("panic: %v", r), WorkerID: id}
		}
	}()

	for job := range jobs {
		err := job.Page.Extract(job.PageNum)
		if err != nil {
			errors <- JobErr{Err: err, WorkerID: id}
		}
		done <- struct{}{}
	}
}
