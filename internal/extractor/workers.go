package main

import (
	"fmt"
	"sync"

	"github.com/dmikhr/pdfjuicer/internal/extractor"
)

// Job contains data about Page processed and current page number
type Job struct {
	page    extractor.Page
	pageNum int
}

// JobErr stores error for workerID if occurs
type JobErr struct {
	err      error
	workerID int
}

// worker process page extraction
func worker(id int, jobs <-chan Job, errors chan<- JobErr, done chan<- struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	// panic recovery
	defer func() {
		if r := recover(); r != nil {
			errors <- JobErr{err: fmt.Errorf("panic: %v", r), workerID: id}
		}
	}()

	for job := range jobs {
		err := job.page.Extract(job.pageNum)
		if err != nil {
			errors <- JobErr{err: err, workerID: id}
		}
		done <- struct{}{}
	}
}
