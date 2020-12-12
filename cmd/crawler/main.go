package main

import (
	"awesomeProject/internal/crawler"
	"awesomeProject/internal/fetcher"
)

var depth = 6
var startURL = "https://clck.ru/9w"

func main() {
	c := crawler.NewCrawler(fetcher.NewFetcher())
	c.Crawl(startURL, depth)
}
