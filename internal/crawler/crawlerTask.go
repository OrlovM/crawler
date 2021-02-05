package crawler

type Task struct {
	Error      error
	Page       *Page
	fetcher    *fetcher
	ParseURLs  bool
	FoundPages PagesSlice
}

func NewCrawlerTask(fetcher *fetcher, page *Page, ParseURLs bool) *Task {
	return &Task{
		fetcher:   fetcher,
		Page:      page,
		ParseURLs: ParseURLs,
	}
}

func (t *Task) Process() {
	t.fetcher.Fetch(t)
}
