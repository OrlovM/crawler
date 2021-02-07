package crawl

import "sort"

type baseWriter struct {
	ch chan *PagesSlice
}

func (b *baseWriter) start(pages chan Page) {
	b.ch = make(chan *PagesSlice)
	var base PagesSlice
	go func() {
		for p := range pages {
			base = append(base, p)
		}
		b.ch <- &base
		close(b.ch)
	}()
}

func (b *baseWriter) getBase() *PagesSlice {
	base := <-b.ch
	sort.Sort(base)
	return base
}
