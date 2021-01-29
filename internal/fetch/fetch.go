package fetch

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	OK       = iota
	Redirect = iota
	NoData   = iota
)

type fetcher struct {
	client http.Client
}

type FetchResult struct {
	StatusCode int
	Location   string
	Body       []byte
	URL        string
	Depth      int
}

type Fetcher interface {
	Fetch(URL string) FetchResult
}

func NewFetcher() *fetcher {
	var client = &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}}
	return &fetcher{*client}
	//TODO Add request timeout
}

func (f *FetchResult) Status() int {
	switch (*f).StatusCode / 100 {
	case 2:
		return OK
	case 3:
		return Redirect
	}
	return NoData
}

func (fetcher fetcher) Fetch(URL string) FetchResult {
	var fResult = FetchResult{0, "", nil, URL, 0}
	resp, err := fetcher.client.Get(URL)
	if err != nil {
		fmt.Println(err, "http.Get failed")
	} else {
		defer resp.Body.Close()
		fResult.StatusCode = resp.StatusCode
		fResult.Location = resp.Header.Get("Location")
		fResult.Body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err, "ioutil.ReadAll failed", resp)
		}
	}
	return fResult
}
