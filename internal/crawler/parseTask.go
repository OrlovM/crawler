package crawler

type ParseTask struct {
	error       error
	FetchResult *FetchResult
	FoundPages  *PagesSlice
}

func NewParseTask(fetchResult *FetchResult) *ParseTask {
	return &ParseTask{FetchResult: fetchResult}
}

func (p *ParseTask) Process() {
	p.FoundPages = Parse(p.FetchResult)
}
