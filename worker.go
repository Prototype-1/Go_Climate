package main

import "sync"

type Task struct {
	City string
}

type Result struct {
	Weather WeatherData
}

type workerPool struct {
	tasks   chan Task
	results chan Result
	workers int
	wg      sync.WaitGroup
}

func NewWorkerPool(workerCount int) *workerPool {
	return &workerPool{
		tasks:   make(chan Task),
		results: make(chan Result),
		workers: workerCount,
	}
}

func (wp *workerPool) Start() {
	for i := 0; i < wp.workers; i++ {
		wp.wg.Add(1) 
		go wp.worker() 
	}
}

func (wp *workerPool) AddTask(task Task) {
	wp.tasks <- task
}

func (wp *workerPool) Results() <-chan Result {
	return wp.results
}

func (wp *workerPool) Stop() {
	close(wp.tasks) 
	wp.wg.Wait()    
	close(wp.results) 
}

func (wp *workerPool) worker() {
	defer wp.wg.Done()
	for task := range wp.tasks {
		weather := fetchWeatherData(task.City)
		wp.results <- Result{Weather: weather}
	}
}

