package main

import (
	"crawler/internal/crawler"
	"flag"
	"fmt"
	"github.com/OrlovM/go-workerpool"
	"net/url"
)

var base crawler.PagesSlice
type URLSlice []*url.URL
type tasksSlice []workerpool.Task
func main() {
	depth := flag.Int("depth", 2, "Depth refers to how far down into a website's page hierarchy crawler crawls")
	startURL := flag.String("url", "https://clck.ru/9w", "URL to start from")
	maxGoroutines := flag.Int("n", 50, "A maximum number of goroutines work at the same time")
	flag.Parse()
	fetcher := crawler.NewFetcher()
	addToBase := make(chan crawler.Page)
	var buffer tasksSlice
	start, err := url.Parse(*startURL)
	if err != nil {
		fmt.Println(err)
		return
	}
	cache := URLSlice{start}
	tasksInWork := 1
	exit := false
	in := make(chan workerpool.Task)
	out := make(chan workerpool.Task)
	pool := workerpool.NewPool(in, out, *maxGoroutines)
	go pool.Run()
	go func() {
		for p := range addToBase {
			base = append(base, p)
		}
	}()
	go func() {
		in <- crawler.NewCrawlerTask(fetcher, &crawler.Page{URL: start, Source: "Set in command line"}, true)
	}()

	for {
		if exit == true {
			break
		}
		select {
		case t := <-out:
			tasksInWork--
			if cT, ok := t.(*crawler.CrawlerTask); ok == true {
				if cT.Error != nil {
					fmt.Println(cT.Error)
					break
				}
				addToBase <- *cT.Page
				printFound(cT.Page)
				for _, p := range cT.FoundPages {
					page := p
					if !cache.contains(page) {
						cache = append(cache, page.URL)
						task := crawler.NewCrawlerTask(fetcher, &page, page.Depth < *depth)
						select {
						case in <- task:
							tasksInWork++
						default:
							buffer = append(buffer, task)
						}
					} else {
						printInBase(&page)
					}
				}
			} //TODO write else branch
		default:
			if len(buffer) != 0 {
				select {
				case in <- buffer[0]:
					tasksInWork++
					buffer.TrimFirst()
				default:
				}
			} else {
				if tasksInWork == 0 {
					exit = true
				}
			}
		}
	}


	fmt.Println("URL in base", len(base))
}

func (s URLSlice) contains(p crawler.Page) bool {
	for _, u := range s {
		if p.URL == u {
			return true
		}
	}
	return false
}

func printInBase(currentPage *crawler.Page) {
	fmt.Println("URL:", currentPage.URL, "Source:", currentPage.Source, " is already in base")
}

func printFound(currentPage *crawler.Page) {
	fmt.Println("Found new URL:", currentPage.URL, "Source:", currentPage.Source, "Depth:", currentPage.Depth, "Status code:", currentPage.StatusCode)
}

func (s *tasksSlice) TrimFirst() {
	*s = (*s)[1:]
}
