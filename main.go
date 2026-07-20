package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

var urls = []string{
	"https://go.dev",
	"https://github.com",
	"https://wikipedia.org",
	"https://news.ycombinator.com",
	"https://reddit.com",
	"https://stackoverflow.com",
}

// numWorkers is pool size: how many goroutines will be processed
// Fixed and controlled - Pattern core.
const numWorkers = 5

// Result load the result after process URL error or not.
type Result struct {
	URL   string
	Title string
	Err   error
}

func main() {
	start := time.Now()

	jobs := make(chan string)
	results := make(chan Result)

	// 1) Sobe o pool de workers.
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(&wg, jobs, results)
	}

	// 2) Feeds jobs into a separate goroutine and closes the channel upon completion.
	go func() {
		for _, url := range urls {
			jobs <- url
		}
		close(jobs)
	}()

	// 3) Close channel after all jobs processed
	go func() {
		wg.Wait()
		close(results)
	}()

	// 4) Collect result until the results channel is closed
	for res := range results {
		if res.Err != nil {
			fmt.Printf("error on %s: %v\n", res.URL, res.Err)
			continue
		}
		fmt.Printf("%s -> %s\n", res.URL, res.Title)
	}

	fmt.Printf("\ndone in %v\n", time.Since(start))
}

// worker consumes URLs from channel jobs, process each one and send Result.
func worker(wg *sync.WaitGroup, jobs <-chan string, results chan<- Result) {
	defer wg.Done()
	for url := range jobs {
		title, err := fetchTitle(url)
		results <- Result{URL: url, Title: title, Err: err}
	}
}

func fetchTitle(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return extractTitle(string(body)), nil
}

func extractTitle(html string) string {
	lower := strings.ToLower(html)
	start := strings.Index(lower, "<title>")
	if start == -1 {
		return "(sem título)"
	}
	start += len("<title>")
	end := strings.Index(lower[start:], "</title>")
	if end == -1 {
		return "(sem título)"
	}
	return strings.TrimSpace(html[start : start+end])
}
