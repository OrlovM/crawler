package crawl

import (
	"net/url"

	"github.com/OrlovM/go-workerpool"
)

type urlSlice []*url.URL
type tasksSlice []workerpool.Task

//Checks if urlSlice contains Page.URL
func (s urlSlice) Contains(p *Page) bool {
	for _, u := range s {
		if *p.URL == *u {
			return true
		}
	}
	return false
}

//Trims the first task in tasksSlice
func (s *tasksSlice) TrimFirst() {
	*s = (*s)[1:]
}
