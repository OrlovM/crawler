package main

import (
	"crawler/internal/crawl"
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	app := cli.App{
		Name:   "Crawler",
		Usage:  "Recursively crawl urls from defined url until the desired depth of recursion will be achieved.",
		Action: startCrawl,
		Authors: []*cli.Author{
			{
				Name:  "Orlov Mikhail",
				Email: "orlov@email.com",
			},
		},
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
}

func startCrawl(ctx *cli.Context) error {
	depth := ctx.Int("depth")
	startURL := ctx.String("startURL")
	concurrency := ctx.Int("concurrency")
	verbose := ctx.Bool("verbose")
	base, err := crawl.Crawl(&startURL, &depth, &concurrency, &verbose)
	if err != nil {
		return err
	}
	fmt.Println("URL in base:", len(*base))
	return nil
}
