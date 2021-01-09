package main

import (
	"crawler/internal/crawler"
	"crawler/internal/fetcher"
)

var depth = 3
var startURL = "https://clck.ru/9w"

func main() {
	c := crawler.NewCrawler(fetcher.NewFetcher())
	c.Crawl(startURL, depth)
}
