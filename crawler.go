package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

const (
	OK       = iota
	Redirect = iota
	NoData   = iota
)

type pageArray []page

var Base pageArray
var PagesToBeScanned pageArray
var depth = 6
var startURL = "https://clck.ru/9w"

var client = &http.Client{
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}}

//TODO Add request timeout

type page struct {
	URL        string
	depth      int
	from       string // URL where this page was found
	statusCode int    //0 if page was not requested or request failed
}

type fetchResult struct {
	statusCode int
	location   string
	body       []byte
}

func (f *fetchResult) status() int {
	switch (*f).statusCode / 100 {
	case 2:
		return OK
	case 3:
		return Redirect
	}
	return NoData

}

func main() {

	PagesToBeScanned = pageArray{page{startURL, 0, "Start URL", 0}}
	crawlRecursive()
	fmt.Println("Unique URLs found", len(Base))

}

func (s *pageArray) trimFirst() {
	*s = (*s)[1:]
}

func crawlRecursive() {

	if len(PagesToBeScanned) != 0 {
		currentPage := PagesToBeScanned[0]
		PagesToBeScanned.trimFirst()
		if Base.contains(currentPage) {
			printInBase(currentPage)
		} else {
			fResult := fetch(currentPage.URL)
			currentPage.statusCode = fResult.statusCode
			Base = append(Base, currentPage)
			printFound(currentPage)
			if fResult.status() == OK {
				if currentPage.depth < depth {
					URLs := extractURLs(fResult.body)
					for _, u := range URLs {
						PagesToBeScanned = append(PagesToBeScanned, page{u, currentPage.depth + 1, currentPage.URL, 0})
					}
				}
			} else if fResult.status() == Redirect {
				if currentPage.depth < depth {
					PagesToBeScanned = append(PagesToBeScanned, page{fResult.location, currentPage.depth, "redirected from" + currentPage.URL, 0})
				}

			} else if fResult.status() == NoData {
				fmt.Println("4xx or 5xx status code recieved on", currentPage.URL)
				//	TODO What to do with this code?
			}
		}
		crawlRecursive()
	}
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

func fetch(URL string) fetchResult {

	var fResult = fetchResult{0, "", nil}
	resp, err := client.Get(URL)
	if err != nil {
		fmt.Println(err, "http.Get failed")
	} else {
		defer resp.Body.Close()
		fResult.statusCode = resp.StatusCode
		fResult.location = resp.Header.Get("Location")
		fResult.body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err, "ioutil.ReadAll failed", resp)
		}
	}
	return fResult

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
