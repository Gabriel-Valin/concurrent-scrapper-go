## Layer 0 ([commit](https://github.com/Gabriel-Valin/concurrent-scrapper-go/commit/2392faa0f89b08e1ae53d57a2ec846abb0963f2f))
- Sequencial scrapper
- Waiting for all request to be completed ("wrong way")
- No HTTP timeout, retries, 0 config config on HTTP CLIENT

## Problem
- We're waiting for all requests to be completed then if one request take too much time like 6~seconds we need to wait this one request finish to keep forward. Also if this request breaks down our program will be generate an error

--------
