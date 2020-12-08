package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)


const (
	OK = iota
	Redirect = iota
	NoData = iota
)

type pageArray []page

var Base pageArray
var PagesToBeScanned pageArray
var depth = 6
var startURL = "https://yandex.ru/ksdfjsdf"

type page struct {
	URL string
	depth int
	from string // URL where this page was found
	statusCode int //0 if page was not requested
}

type fetchResult struct {
	statusCode int
	location   string
	body       []byte
}

func (f fetchResult) status() int {
	switch f.statusCode / 100 {
	case 2 : return  OK
	case 3 : return  Redirect
	}
	return NoData

}

func main() {

	PagesToBeScanned = pageArray{page{startURL,0, "Start URL", 0}}
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
				//	TODO Handle redirect

			} else if fResult.status() == NoData {
				fmt.Println("NODATA!!!!!!!!!!!!!!")
				//	TODO What to do with this code?
			} else {
				//	TODO Something went wrong
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
	fmt.Println("Found new URL", currentPage.URL,"on page", currentPage.from, "depth", currentPage.depth, "status code", currentPage.statusCode)
}

func fetch(URL string) fetchResult {


	var fResult = fetchResult{0, "", nil}


	resp, err := http.Get(URL)
	if err != nil {
		fmt.Println(err, "http.Get не получилось")
	} else {

		defer resp.Body.Close()

		fResult.body, err = ioutil.ReadAll(resp.Body)

		fResult.statusCode = resp.StatusCode

		if err != nil {
			fmt.Println(err, "ioutil.ReadAll не получилось", resp)
		}
	}
	return fResult

}

func extractURLs(body []byte) []string {

	var URLsFound []string
	re := regexp.MustCompile(`href="http.+?"`)
	strings := re.FindAllString(string(body), -1)
	for _, v := range strings {
		currentURL := v[len("href='"):(len(v)-1)]
		URLsFound = append(URLsFound, currentURL, currentURL)
	}
	return URLsFound
}



