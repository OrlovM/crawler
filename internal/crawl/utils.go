package crawl

import (
	"github.com/OrlovM/go-workerpool"
	"net/url"
)

type URLSlice []*url.URL
type TasksSlice []workerpool.Task

func (s URLSlice) Contains(p *Page) bool {
	for _, u := range s {
		if *p.URL == *u {
			return true
		}
	}
	return false
}

func (s *TasksSlice) TrimFirst() {
	*s = (*s)[1:]
}
