package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Job struct {
	URL string
}

type Result struct {
	URL    string
	Status int
}

func worker(id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := range jobs {
		fmt.Printf("Worker %d: Checking %s\n", id, i.URL)
		// Actual HTTP CALL
		resp, err := http.Get(i.URL)

		status := 0
		if err == nil {
			status = resp.StatusCode
			resp.Body.Close()
		}

		results <- Result{
			URL:    i.URL,
			Status: status,
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
	jobs := make(chan Job, 5)
	results := make(chan Result, 5)

	//workers creation
	for w := 1; w <= 3; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	//URLS
	urls := []string{"https://www.youtube.com", "https://www.google.com", "https://www.linkedin.com", "https://www.naukri.com"}

	// jobs creation
	for i := range urls {
		url := Job{
			URL: urls[i],
		}
		jobs <- url
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	//captures result
	for res := range results {
		fmt.Printf("Result: %s -> Status: %d\n", res.URL, res.Status)
	}

	fmt.Printf("\nDone! Total time taken: %s\n", time.Since(start))

	//fmt.Println("Hello, world")
}
