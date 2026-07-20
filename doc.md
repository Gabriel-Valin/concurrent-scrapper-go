## Layer 0
- Sequencial scrapper
- Waiting for all request to be completed ("wrong way")
- No HTTP timeout, retries, 0 config config on HTTP CLIENT

## Problem
- We're waiting for all requests to be completed then if one request take too much time like 6~seconds we need to wait this one request finish to keep forward. Also if this request breaks down our program will be generate an error
