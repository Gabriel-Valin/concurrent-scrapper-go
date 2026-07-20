## Layer 0 ([commit](https://github.com/Gabriel-Valin/concurrent-scrapper-go/commit/2392faa0f89b08e1ae53d57a2ec846abb0963f2f))
- Sequencial scrapper
- Waiting for all request to be completed ("wrong way")
- No HTTP timeout, retries, 0 config config on HTTP CLIENT

## Problem
- We're waiting for all requests to be completed then if one request take too much time like 6~seconds we need to wait this one request finish to keep forward. Also if this request breaks down our program will be generate an error

--------

## Layer 1 ([commit](https://github.com/Gabriel-Valin/concurrent-scrapper-go/commit/384f7ea5d7fefc78f851a2b3f0c5804a67291272))
- Using worker pool with solid work numbers to processed each URL
- Using sync.WaitGroup and channels to share communication between goroutines and control workflow properly
- Now we're processing each URL concurrently and the most slow request will be our "bottleneck"

## Problem
- That `http.Get` without a timeout I criticized in Layer 0? It’s a ticking time bomb now. If one of those sites hangs and never responds, the worker that picked it up gets stuck in that `fetchTitle` call forever. It never returns to the jobs range, never processes another URL—you’ve permanently lost a worker. If several of them hang, the entire pool freezes, `wg.Wait()` never completes, the `results` channel never closes, and the main routine hangs. A single slow site brings the whole thing down.
