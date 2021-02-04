package crawler

type CrawlerTask struct {
	Error      error
	Page       *Page
	fetcher    *fetcher
	ParseURLs  bool
	FoundPages PagesSlice
}

func NewCrawlerTask(fetcher *fetcher, page *Page, ParseURLs bool) *CrawlerTask {
	return &CrawlerTask{
		fetcher:   fetcher,
		Page:      page,
		ParseURLs: ParseURLs,
	}
}

func (t *CrawlerTask) Process() {
	t.fetcher.Fetch(t)
}
