package main

import (
	"fmt"
	"sync"
)

type Job struct {
	page    Page
	pageNum int
}

type JobErr struct {
	err      error
	workerID int
}

func worker(id int, jobs <-chan Job, errors chan<- JobErr, done chan<- struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	defer func() {
		if r := recover(); r != nil {
			errors <- JobErr{err: fmt.Errorf("panic: %v", r), workerID: id}
		}
	}()

	for job := range jobs {
		err := job.page.extract(job.pageNum)
		if err != nil {
			errors <- JobErr{err: err, workerID: id}
		}
		done <- struct{}{}
	}
}
