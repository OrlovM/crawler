package crawler

import (
	"crawler/internal/fetcher"
	"fmt"
	"regexp"
	"time"
)

type page struct {
	URL        string
	depth      int
	from       string // URL where this page was found
	statusCode int    //0 if page was not requested or request failed
}

type pageArray []page

type crawler struct {
	fetcher fetcher.Fetcher
}

var Base pageArray
var depth int

func NewCrawler(fetcher fetcher.Fetcher) *crawler {
	return &crawler{fetcher}
}

func (c *crawler) Crawl(startURL string, maxDepth int, maxGoroutines int) {
	depth = maxDepth
	sem := make(chan struct{}, maxGoroutines)
	pagesFound := make(chan *page, 10000)
	pagesToCrawl := make(chan *page, 10000)
	pagesFound <- &page{startURL, 0, "Start URL", 0}
	go showInfo(pagesFound, pagesToCrawl)
	go operate(pagesFound, pagesToCrawl)
	for p := range pagesToCrawl {
		sem <- struct{}{}
		go c.crawl(p, pagesFound, sem)
	}

	fmt.Println("Unique URLs found", len(Base))
}

func showInfo(pagesFound chan *page, pagesToCrawl chan *page) {
	time.Sleep(time.Second * 1)
	fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!pages Found", len(pagesFound), "pages to crawl", len(pagesToCrawl), "base", len(Base))
	showInfo(pagesFound, pagesToCrawl)
}

func operate(pagesFound chan *page, pagesToCrawl chan *page) {

	for p := range pagesFound {
		if Base.contains(p) {
			printInBase(*p)
		} else {
			printFound(*p)
			Base = append(Base, *p)
			if p.depth < depth {
				pagesToCrawl <- p
			}
		}
	}
}

func (c *crawler) crawl(currentPage *page, pagesFound chan *page, sem chan struct{}) {
	defer func() { <-sem }()
	fResult := c.fetcher.Fetch(currentPage.URL)
	currentPage.statusCode = fResult.StatusCode
	if fResult.Status() == fetcher.OK {
		URLs := extractURLs(fResult.Body)
		for _, u := range URLs {
			pagesFound <- &page{u, currentPage.depth + 1, currentPage.URL, 0}
		}
	} else if fResult.Status() == fetcher.Redirect {
		pagesFound <- &page{fResult.Location, currentPage.depth, "redirected from" + currentPage.URL, 0}
	} else if fResult.Status() == fetcher.NoData {
		fmt.Println("4xx or 5xx status code recieved on", currentPage.URL)
		//	TODO What to do with this code?
	}

}

func (s *pageArray) trimFirst() {
	*s = (*s)[1:]
}

func (s pageArray) contains(p *page) bool {
	for _, n := range s {
		if p.URL == n.URL {
			return true
		}
	}
	return false
}

func printInBase(currentPage page) {
	fmt.Println("URL ", currentPage.URL, "found on page", currentPage.from, " is already in base")
}

func printFound(currentPage page) {
	fmt.Println("Found new URL", currentPage.URL, "on page", currentPage.from, "depth", currentPage.depth, "status code", currentPage.statusCode)
}

func extractURLs(body []byte) []string {
	var URLsFound []string
	re := regexp.MustCompile(`href="http.+?"`)
	// TODO Find a better way to parse URLs, handle local links
	strings := re.FindAllString(string(body), -1)
	for _, v := range strings {
		currentURL := v[len("href='"):(len(v) - 1)]
		URLsFound = append(URLsFound, currentURL, currentURL)
	}
	return URLsFound
}
