package main

import (
	"crawler/internal/crawler"
	"crawler/internal/fetch"
	"flag"
	"fmt"
	"sync"
	"time"
)

var base crawler.PageArray
var Start = time.Now()

type URLSlice []string

func main() {
	depth := flag.Int("depth", 2, "Depth refers to how far down into a website's page hierarchy crawler crawls")
	startURL := flag.String("url", "https://ya.ru", "URL to start from")
	maxGoroutines := flag.Int("n", 50, "A maximum number of goroutines work at the same time")
	flag.Parse()
	fetcher := fetch.NewFetcher()
	unfetchedPages := make(chan crawler.Page)
	addToBase := make(chan crawler.Page)
	foundPagesChannels := make(chan chan crawler.Page)
	fetchChannels := make(chan chan fetch.FetchResult)
	quit := make(chan struct{})

	Start = time.Now()

	//base manager
	go func() {
		for p := range addToBase {
			base = append(base, p)
		}
	}()

	//crating manager
	go func() {
		foundPages := mergeParsed(foundPagesChannels)
		foundPages <- crawler.Page{URL: *startURL}
		var cache URLSlice
		var buffer crawler.PageArray
		for {
			select {
			case p, ok := <-foundPages:
				if ok {
					if !cache.contains(p) {
						fmt.Println("found new url", p.URL)
						cache = append(cache, p.URL)
						select {
						case unfetchedPages <- p:
						default:
							buffer = append(buffer, p)
						}
					} else {
						println("URL", p.URL, "is already in base")
					}
				} else {
					return
				}
			case <-quit:
				close(unfetchedPages)
				return
			default:
				if len(buffer) != 0 {
					select {
					case unfetchedPages <- buffer[0]:
						buffer.TrimFirst()
					default:
						break
					}
				}
			}
		}
	}()

	// creating parsers pool
	go func() {
		fetchResults := mergeFetched(fetchChannels)
		for i := 0; i < (1 + *maxGoroutines/5); i++ {
			foundPages := make(chan crawler.Page)
			foundPagesChannels <- foundPages
			go func() {
				for f := range fetchResults {
					for _, u := range crawler.ExtractURLs(f.Body) {
						foundPages <- crawler.Page{URL: u, Depth: f.Depth + 1, From: f.URL}
					}
				}
				close(foundPages)
			}()
		}
		close(foundPagesChannels)
	}()

	// creating fetchers pool
	for i := 0; i < *maxGoroutines; i++ {
		fResults := make(chan fetch.FetchResult)
		fetchChannels <- fResults
		go func() {
			for p := range unfetchedPages {
				fResult := fetcher.Fetch(p.URL) //TODO add error check
				p.StatusCode = fResult.StatusCode
				addToBase <- p
				fmt.Println("added to base", p.URL)
				if p.Depth < *depth { //TODO handle redirect
					fResult.Depth = p.Depth
					fResults <- fResult
				}
			}
			close(fResults)
		}()

	}
	close(fetchChannels)

	time.Sleep(100000 * time.Millisecond)

	fmt.Println("URL in base", len(base))

}

//Merge fetchers output channels in one
func mergeFetched(fetchChannels chan chan fetch.FetchResult) <-chan fetch.FetchResult {
	out := make(chan fetch.FetchResult, 20)
	var wg sync.WaitGroup

	output := func(c <-chan fetch.FetchResult) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}

	for c := range fetchChannels {
		wg.Add(1)
		go output(c)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

//Merge fetchers output channels in one
func mergeParsed(foundPagesChannels chan chan crawler.Page) chan crawler.Page {

	out := make(chan crawler.Page, 20)
	var wg sync.WaitGroup

	output := func(c <-chan crawler.Page) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}

	for c := range foundPagesChannels {
		wg.Add(1)
		go output(c)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func (s URLSlice) contains(p crawler.Page) bool {
	for _, n := range s {
		if p.URL == n {
			return true
		}
	}
	return false
}
