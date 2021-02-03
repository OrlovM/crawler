package main

import (
	"crawler/internal/crawler"
	"flag"
	"fmt"
	"github.com/OrlovM/go-workerpool"
	"time"
)

var base crawler.PagesSlice
var Start = time.Now()

type URLSlice []string

func main() {
	depth := flag.Int("depth", 3, "Depth refers to how far down into a website's page hierarchy crawler crawls")
	startURL := flag.String("url", "https://ya.ru", "URL to start from")
	maxGoroutines := flag.Int("n", 50, "A maximum number of goroutines work at the same time")
	flag.Parse()
	fetcher := crawler.NewFetcher()
	fIn := make(chan workerpool.Task)
	fOut := make(chan workerpool.Task)
	pIn := make(chan workerpool.Task)
	pOut := make(chan workerpool.Task)
	addToBase := make(chan crawler.Page)
	var buffer []workerpool.Task
	cache := URLSlice{*startURL}
	tasksInWork := 1
	exit := false
	fetchersPool := workerpool.NewPool(fIn, fOut, *maxGoroutines)
	parsersPool := workerpool.NewPool(pIn, pOut, 1+*maxGoroutines/5)
	go fetchersPool.Run()
	go parsersPool.Run()
	go func() {
		for p := range addToBase {
			base = append(base, p)
		}
	}()
	go func() {
		fIn <- crawler.NewFetchTask(fetcher, &crawler.Page{URL: *startURL, Depth: 0})
	}()
	for {
		if exit == true {break}
		select {
		case t := <- fOut:
			tasksInWork--
			if fT, ok := t.(*crawler.FetchTask); ok == true {
				addToBase <- *fT.Page
				if fT.Page.Depth < *depth {
				pTask := crawler.NewParseTask(fT.FetchResult)
				select {
				case pIn <- pTask:
					tasksInWork++
				default:
					buffer = append(buffer, pTask)
					}
				}
			} //TODO write else branch
		case t := <- pOut:
			tasksInWork--
			if pT, ok := t.(*crawler.ParseTask); ok == true {
				for _, p := range *pT.FoundPages {
					page := p
					fmt.Println(p.URL)
					if !cache.contains(page) {
						printFound(page)
						cache = append(cache, page.URL)
						fTask := crawler.NewFetchTask(fetcher, &page)
						select {
						case fIn <- fTask:
							tasksInWork++
						default:
							buffer = append(buffer, fTask)
						}
					} else {printInBase(page)}
				}
			}
		default:
			//fmt.Println("default")
			if len(buffer) != 0 {
				switch t := buffer[0].(type) {
				case *crawler.FetchTask:
					select {
					case fIn <- t:
						tasksInWork++
						buffer = buffer[1:]
					default:
					}
				case *crawler.ParseTask:
					select {
					case pIn <- t:
						tasksInWork++
					default:
					}
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
	for _, n := range s {
		if p.URL == n {
			return true
		}
	}
	return false
}

func printInBase(currentPage crawler.Page) {
	fmt.Println("URL ", currentPage.URL, "found on Page", currentPage.From, " is already in base")
}

func printFound(currentPage crawler.Page) {
	fmt.Println("Found new URL", currentPage.URL, "on Page", currentPage.From, "Depth", currentPage.Depth, "status code", currentPage.StatusCode)
}
