package crawler

type FetchTask struct {
	error       error
	Page        *Page
	FetchResult *FetchResult
	fetcher     *fetcher
}

func NewFetchTask(fetcher *fetcher, page *Page) *FetchTask {
	return &FetchTask{
		fetcher: fetcher,
		Page:    page,
	}
}

func (f *FetchTask) Process() {
	f.FetchResult = f.fetcher.Fetch(f.Page.URL)
	f.FetchResult.Depth = f.Page.Depth
	f.Page.StatusCode = f.FetchResult.StatusCode
}
