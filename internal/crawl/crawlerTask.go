package crawl

type task struct {
	Error      error
	Page       *Page
	fetcher    *fetcher
	ParseURLs  bool
	FoundPages PagesSlice
}

//newCrawlerTask creates and returns task
func newCrawlerTask(fetcher *fetcher, page *Page, ParseURLs bool) *task {
	return &task{
		fetcher:   fetcher,
		Page:      page,
		ParseURLs: ParseURLs,
	}
}

func (t *task) Process() {
	t.fetcher.Fetch(t)
}
