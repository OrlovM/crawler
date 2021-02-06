package crawl

import (
	"fmt"
	"github.com/OrlovM/go-workerpool"
	"net/url"
	"sync"
)

func Crawl(startURL *string, depth *int, concurrency *int, verbose *bool) (*PagesSlice, error) {
	var (
		errors []error
		base   PagesSlice
		buffer TasksSlice
		wg     sync.WaitGroup
	)
	printer := NewPrinter(*verbose)
	fetcher := NewFetcher()
	addToBase := make(chan Page)
	in := make(chan workerpool.Task)
	out := make(chan workerpool.Task)
	start, err := url.Parse(*startURL)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	cache := URLSlice{start}
	tasksInWork := 1
	exit := false
	pool := workerpool.NewPool(in, out, *concurrency)
	go pool.Run()
	wg.Add(1)
	go func() {
		for p := range addToBase {
			base = append(base, p)
		}
		wg.Done()
	}()
	go func() {
		in <- NewCrawlerTask(fetcher, &Page{URL: start, Source: "Set in command line"}, true)
	}()

	for {
		if exit == true {
			break
		}
		select {
		case t := <-out:
			tasksInWork--
			if cT, ok := t.(*Task); ok == true {
				if cT.Error != nil {
					printer.Error(cT.Error)
					errors = append(errors, cT.Error)
					break
				}
				addToBase <- *cT.Page
				for _, p := range cT.FoundPages {
					page := p
					if !cache.Contains(&page) {
						printer.Found(&page)
						cache = append(cache, page.URL)
						task := NewCrawlerTask(fetcher, &page, page.Depth < *depth)
						select {
						case in <- task:
							tasksInWork++
						default:
							buffer = append(buffer, task)
						}
					} else {
						printer.InBase(&page)
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
					close(addToBase)
					close(in)
					exit = true
				}
			}
		}
	}
	wg.Wait()
	return &base, nil
}
