package main

import (
	"crawler/internal/crawl"
	"crawler/internal/tools"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/urfave/cli/v2"
)

type Result struct {
	StartURL        string           `json:"start_url"`
	UniqueURLsFound int              `json:"unique_ur_ls_found"`
	ErrorsOccured   int              `json:"errors_occured"`
	URLs            crawl.PagesSlice `json:"urls"`
	Errors          []error          `json:"errors"`
}

func main() {
	app := cli.App{
		Name:   "Crawler",
		Usage:  "Recursively crawl urls from defined url until the desired depth of recursion will be achieved.",
		Action: Run,
		Authors: []*cli.Author{
			{
				Name:  "Orlov Mikhail",
				Email: "orlov@email.com",
			},
		},
		Compiled: time.Now(),
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
			&cli.StringFlag{
				Name:  "filepath",
				Value: "crawl_results/",
				Usage: `A path there the file with results of crawl will be created.
				Could be specified with file name e.g. "/a/b.txt" or only a directory "/a/".
				If file name is not specified it will be set to default pattern "2006-01-02T15:04:05.json"`,
				Aliases: []string{"p"}},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func Run(ctx *cli.Context) error {
	results, err := tools.CreateFile(ctx)
	if err != nil {
		log.Fatal("Unable to create file:", err)
	}

	depth := ctx.Int("depth")
	startURL := ctx.String("startURL")
	concurrency := ctx.Int("concurrency")
	verbose := ctx.Bool("verbose")

	base, errors, err := crawl.Crawl(&startURL, &depth, &concurrency, &verbose)
	if err != nil {
		return err
	}

	r := Result{string(startURL), len(*base), len(*errors), *base, *errors}
	j, err := json.MarshalIndent(r, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	_, err = results.Write(j)
	if err != nil {
		log.Fatal(err)
	}

	err = results.Close()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("URL in base:", len(*base))
	fmt.Println("Errors occurred", len(*errors))

	return nil
}
