package main

import (
	"crawler/internal/crawler"
	"crawler/internal/fetcher"
	"flag"
)

func main() {

	depth := flag.Int("depth", 5, "Depth refers to how far down into a website's page hierarchy crawler crawls")
	startURL := flag.String("url", "https://clck.ru/9w", "URL to start from")
	maxGoroutines := flag.Int("n", 20, "A maximum number of goroutines work at the same time")

	c := crawler.NewCrawler(fetcher.NewFetcher())
	c.Crawl(*startURL, *depth, *maxGoroutines)
}
