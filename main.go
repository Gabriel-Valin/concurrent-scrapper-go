package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os/signal"
	"strings"
	"sync"
	"syscall"
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

const numWorkers = 5

// requestTimeout é o prazo máximo de UMA requisição. Passou disso, aborta aquela.
const requestTimeout = 10 * time.Second

type Result struct {
	URL   string
	Title string
	Err   error
}

func main() {
	start := time.Now()

	// Context raiz: cancelado por Ctrl+C (SIGINT) ou SIGTERM.
	// signal.NotifyContext devolve um ctx que "morre" quando o sinal chega.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	jobs := make(chan string)
	results := make(chan Result)

	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(ctx, &wg, jobs, results)
	}

	// Alimenta os jobs, mas respeitando o cancelamento: se ctx morrer,
	// paramos de enfileirar em vez de insistir.
	go func() {
		defer close(jobs)
		for _, url := range urls {
			select {
			case jobs <- url:
			case <-ctx.Done():
				return
			}
		}
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	for res := range results {
		if res.Err != nil {
			fmt.Printf("erro em %s: %v\n", res.URL, res.Err)
			continue
		}
		fmt.Printf("%s -> %s\n", res.URL, res.Title)
	}

	if ctx.Err() != nil {
		fmt.Printf("\ninterrompido: %v\n", ctx.Err())
	}
	fmt.Printf("concluído em %v\n", time.Since(start))
}

func worker(ctx context.Context, wg *sync.WaitGroup, jobs <-chan string, results chan<- Result) {
	defer wg.Done()
	for url := range jobs {
		title, err := fetchTitle(ctx, url)
		// Envia o resultado, mas sem travar se o programa está encerrando.
		select {
		case results <- Result{URL: url, Title: title, Err: err}:
		case <-ctx.Done():
			return
		}
	}
}

// fetchTitle agora recebe context e impõe um timeout por requisição.
func fetchTitle(ctx context.Context, url string) (string, error) {
	// Deriva um context com prazo: o menor entre o timeout da requisição
	// e o cancelamento global herdado de ctx.
	ctx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
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
