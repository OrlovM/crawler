package crawl

import "net/url"

//Page is a representation of web page
type Page struct {
	URL        *url.URL
	Depth      int
	Source     string // URL where this Page was found
	StatusCode int    //0 if Page was not requested or request failed
}

//PagesSlice is a slice of Pages. Made for ability to implement methods on slice.
type PagesSlice []Page

//Implementation of sort interface
func (s PagesSlice) Len() int { return len(s) }

func (s PagesSlice) Less(i, j int) bool { return s[i].Depth < s[j].Depth }

func (s PagesSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
