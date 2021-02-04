package crawler

import "net/url"

type Page struct {
	URL        *url.URL
	Depth      int
	Source     string // URL where this Page was found
	StatusCode int    //0 if Page was not requested or request failed
}

type PagesSlice []Page

func (s PagesSlice) contains(p *Page) bool {
	for _, n := range s {
		if p.URL == n.URL {
			return true
		}
	}
	return false
}
