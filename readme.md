# Crawler

A simple web crawler written on Go. It recursively crawls urls from defined url untill the desired depth of recursion will be achieved.

## How to use

Download and build project.

Run it with or without command line flags.

Example:

```
./crawler -depth 5 -n 20 -url https://google.com
```

Command line flags:

-depth Depth refers to how far down into a website's page hierarchy crawler crawls (default 5)

-n A maximum number of goroutines work at the same time (default 20)

-url URL to start from (default "https://clck.ru/9w")

