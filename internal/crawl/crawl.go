package crawl

import (
	"github.com/OrlovM/go-workerpool"
	"net/url"
)

func Crawl(startURL *string, depth *int, concurrency *int, verbose *bool) (*PagesSlice, error) {
	var (
		errors      []error
		cache       URLSlice
		tasksInWork int
		exit        bool

		addToBase = make(chan Page)
		in        = make(chan workerpool.Task)
		out       = make(chan workerpool.Task)

		bw      = baseWriter{}
		printer = NewPrinter(*verbose)
		fetcher = NewFetcher()
		pool    = workerpool.NewPool(in, out, *concurrency)
	)

	start, err := url.Parse(*startURL)
	if err != nil {
		return nil, err
	}
	buffer := TasksSlice{NewCrawlerTask(fetcher, &Page{URL: start, Source: "Set in command line"}, true)}

	go pool.Run()
	bw.start(addToBase)

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
	return bw.getBase(), nil
}
