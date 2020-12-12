package crawler

import (

	//"fmt"

	"awesomeProject/internal/fetcher"
	"fmt"
	"regexp"
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
var PagesToBeScanned pageArray
var depth int

func NewCrawler(fetcher fetcher.Fetcher) *crawler {
	return &crawler{fetcher}
}

func (c *crawler) Crawl(startURL string, maxDepth int) {
	PagesToBeScanned = pageArray{page{startURL, 0, "Start URL", 0}}
	depth = maxDepth
	c.crawlRecursive()
	fmt.Println("Unique URLs found", len(Base))
}

func (c *crawler) crawlRecursive() {
	if len(PagesToBeScanned) != 0 {
		currentPage := PagesToBeScanned[0]
		PagesToBeScanned.trimFirst()
		if Base.contains(currentPage) {
			printInBase(currentPage)
		} else {
			fResult := c.fetcher.Fetch(currentPage.URL)
			currentPage.statusCode = fResult.StatusCode
			Base = append(Base, currentPage)
			printFound(currentPage)
			if fResult.Status() == fetcher.OK {
				if currentPage.depth < depth {
					URLs := extractURLs(fResult.Body)
					for _, u := range URLs {
						PagesToBeScanned = append(PagesToBeScanned, page{u, currentPage.depth + 1, currentPage.URL, 0})
					}
				}
			} else if fResult.Status() == fetcher.Redirect {
				if currentPage.depth < depth {
					PagesToBeScanned = append(PagesToBeScanned, page{fResult.Location, currentPage.depth, "redirected from" + currentPage.URL, 0})
				}

			} else if fResult.Status() == fetcher.NoData {
				fmt.Println("4xx or 5xx status code recieved on", currentPage.URL)
				//	TODO What to do with this code?
			}
		}
		c.crawlRecursive()
	}
}

func (s *pageArray) trimFirst() {
	*s = (*s)[1:]
}

func (s pageArray) contains(p page) bool {
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
