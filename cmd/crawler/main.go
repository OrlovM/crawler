package main

import (
	"crawler/internal/crawler"
	"crawler/tools"
	"crawler/tools/print"
	"fmt"
	"github.com/OrlovM/go-workerpool"
	"github.com/urfave/cli/v2"
	"log"
	"net/url"
	"os"
	"sync"
)

var base crawler.PagesSlice

func main() {

	app := cli.App{
		Name:   "Crawler",
		Usage:  "Recursively crawl urls from defined url until the desired depth of recursion will be achieved.",
		Action: crawl,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "startURL",
				Value: "https://ya.ru",
				Usage: "URL to start from"},
			&cli.IntFlag{
				Name:  "depth",
				Value: 2,
				Usage: "Depth refers to how far down into a website's page hierarchy crawler crawls"},
			&cli.IntFlag{
				Name:  "concurrency",
				Value: 50,
				Usage: "A maximum number of goroutines work at the same time"},
			&cli.BoolFlag{
				Name:    "verbose",
				Value:   false,
				Usage:   "Prints details about crawling process",
				Aliases: []string{"v"}},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("URL in base:", len(base))
}

func crawl(ctx *cli.Context) error {
	depth := ctx.Int("depth")
	startURL := ctx.String("startURL")
	concurrency := ctx.Int("concurrency")
	printer := print.NewPrinter(ctx.Bool("verbose"))
	var errors []error
	fetcher := crawler.NewFetcher()
	addToBase := make(chan crawler.Page)
	var buffer tools.TasksSlice
	start, err := url.Parse(startURL)
	if err != nil {
		fmt.Println(err)
		return err
	}
	cache := tools.URLSlice{start}
	tasksInWork := 1
	exit := false
	in := make(chan workerpool.Task)
	out := make(chan workerpool.Task)
	pool := workerpool.NewPool(in, out, concurrency)
	var wg sync.WaitGroup
	go pool.Run()
	wg.Add(1)
	go func() {
		for p := range addToBase {
			base = append(base, p)
		}
		wg.Done()
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
			if cT, ok := t.(*crawler.Task); ok == true {
				if cT.Error != nil {
					printer.Error(cT.Error)
					errors = append(errors, cT.Error)
					break
				}
				addToBase <- *cT.Page
				printer.Found(cT.Page)
				for _, p := range cT.FoundPages {
					page := p
					if !cache.Contains(&page) {
						cache = append(cache, page.URL)
						task := crawler.NewCrawlerTask(fetcher, &page, page.Depth < depth)
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
					exit = true
				}
			}
		}
	}
	wg.Wait()
	return nil
}
