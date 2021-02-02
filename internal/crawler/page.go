package crawler

type Page struct {
	URL        string
	Depth      int
	From       string // URL where this Page was found
	StatusCode int    //0 if Page was not requested or request failed
}

type PagesSlice []Page

func (s *PagesSlice) TrimFirst() {
	*s = (*s)[1:]
}

func (s PagesSlice) contains(p *Page) bool {
	for _, n := range s {
		if p.URL == n.URL {
			return true
		}
	}
	return false
}
