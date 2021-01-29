package crawler

import (
	"crawler/internal/fetch"
	"fmt"
	"regexp"
	"time"
)

type Page struct {
	URL        string
	Depth      int
	From       string // URL where this Page was found
	StatusCode int    //0 if Page was not requested or request failed
}

type PageArray []Page

type crawler struct {
	fetcher fetch.Fetcher
}

var Base PageArray
var depth int

func NewCrawler(fetcher fetch.Fetcher) *crawler {
	return &crawler{fetcher}
}

func (c *crawler) Crawl(startURL string, maxDepth int, maxGoroutines int) {
	depth = maxDepth
	sem := make(chan struct{}, maxGoroutines)
	pagesFound := make(chan *Page, 10000)
	pagesToCrawl := make(chan *Page, 10000)
	pagesFound <- &Page{startURL, 0, "Start URL", 0}
	go showInfo(pagesFound, pagesToCrawl)
	go operate(pagesFound, pagesToCrawl)
	for p := range pagesToCrawl {
		sem <- struct{}{}
		go c.crawl(p, pagesFound, sem)
	}

	fmt.Println("Unique URLs found", len(Base))
}

func showInfo(pagesFound chan *Page, pagesToCrawl chan *Page) {
	time.Sleep(time.Second * 1)
	fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!pages Found", len(pagesFound), "pages to crawl", len(pagesToCrawl), "base", len(Base))
	showInfo(pagesFound, pagesToCrawl)
}

func operate(pagesFound chan *Page, pagesToCrawl chan *Page) {

	for p := range pagesFound {
		if Base.contains(p) {
			printInBase(*p)
		} else {
			printFound(*p)
			Base = append(Base, *p)
			if p.Depth < depth {
				pagesToCrawl <- p
			}
		}
	}
}

func (c *crawler) crawl(currentPage *Page, pagesFound chan *Page, sem chan struct{}) {
	defer func() { <-sem }()
	fResult := c.fetcher.Fetch(currentPage.URL)
	currentPage.StatusCode = fResult.StatusCode
	if fResult.Status() == fetch.OK {
		URLs := ExtractURLs(fResult.Body)
		for _, u := range URLs {
			pagesFound <- &Page{u, currentPage.Depth + 1, currentPage.URL, 0}
		}
	} else if fResult.Status() == fetch.Redirect {
		pagesFound <- &Page{fResult.Location, currentPage.Depth, "redirected From" + currentPage.URL, 0}
	} else if fResult.Status() == fetch.NoData {
		fmt.Println("4xx or 5xx status code recieved on", currentPage.URL)
		//	TODO What to do with this code?
	}

}

func (s *PageArray) TrimFirst() {
	*s = (*s)[1:]
}

func (s PageArray) contains(p *Page) bool {
	for _, n := range s {
		if p.URL == n.URL {
			return true
		}
	}
	return false
}

func printInBase(currentPage Page) {
	fmt.Println("URL ", currentPage.URL, "found on Page", currentPage.From, " is already in base")
}

func printFound(currentPage Page) {
	fmt.Println("Found new URL", currentPage.URL, "on Page", currentPage.From, "Depth", currentPage.Depth, "status code", currentPage.StatusCode)
}

func ExtractURLs(body []byte) []string {
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
