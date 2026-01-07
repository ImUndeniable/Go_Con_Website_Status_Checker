package main

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

type Job struct {
	URL string
}

type Result struct {
	URL    string
	Status int
	Error  error
}

func worker(id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := range jobs {
		//Context - 2 secs timeout
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel() // always cancel the free resources
		//create the request with context
		req, err := http.NewRequestWithContext(ctx, "GET", i.URL, nil)
		if err != nil {
			results <- Result{URL: i.URL, Status: 0, Error: err}
			continue
		}

		fmt.Printf("Worker %d: Checking %s\n", id, i.URL)
		// Execute the request
		resp, err := http.DefaultClient.Do(req)

		status := 0
		if err == nil {
			status = resp.StatusCode
			resp.Body.Close()
		}

		results <- Result{
			URL:    i.URL,
			Status: status,
			Error:  err,
		}
		/*final := Result{
			URL:    i.URL,
			Status: 200,
		}
		results <- final */
	}
}

func main() {
	start := time.Now()

	var wg sync.WaitGroup
	jobs := make(chan Job, 35)
	results := make(chan Result, 35)

	//workers creation
	for w := 1; w <= 10; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	file, err := os.Open("websites.txt")
	if err != nil {
		fmt.Println("The file is not readable/avaiable")
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// jobs creation
	go func() {
		for scanner.Scan() {
			url := scanner.Text()
			jobs <- Job{URL: url}
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	//captures result
	for res := range results {
		if res.Error != nil {
			fmt.Printf("Result: %s ->Status: ERROR (%v)\n", res.URL, res.Error)
			continue
		}
		fmt.Printf("Result: %s -> Status: %d\n", res.URL, res.Status)
	}

	fmt.Printf("\nDone! Total time taken: %s\n", time.Since(start))

	//fmt.Println("Hello, world")
}
