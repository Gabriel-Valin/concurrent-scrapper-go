package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
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

func main() {
	start := time.Now()

	for _, url := range urls {
		title, err := fetchTitle(url)
		if err != nil {
			fmt.Printf("error on %s: %v\n", url, err)
			continue
		}
		fmt.Printf("%s -> %s\n", url, title)
	}

	fmt.Printf("\n done in %v\n", time.Since(start))
}

// fetchTitle download page and extract title
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

// extractTitle find content between <title> and </title>
func extractTitle(html string) string {
	lower := strings.ToLower(html)
	start := strings.Index(lower, "<title>")
	if start == -1 {
		return "(no title)"
	}
	start += len("<title>")
	end := strings.Index(lower[start:], "</title>")
	if end == -1 {
		return "(no title)"
	}
	return strings.TrimSpace(html[start : start+end])
}
