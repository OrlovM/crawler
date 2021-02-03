package crawler

import (
	"regexp"
)

func Parse(f *FetchResult) *PagesSlice {
	pages := PagesSlice{}
	for _, u := range ExtractURLs(f.Body) {
		pages = append(pages, Page{u, f.Depth + 1, f.URL, 0})
	}
	return &pages
}

func ExtractURLs(body []byte) []string {
	var URLsFound []string
	re := regexp.MustCompile(`href="http.+?"`)
	// TODO Find a better way to parse URLs, handle local links
	strings := re.FindAllString(string(body), -1)
	for _, v := range strings {
		currentURL := v[len("href='"):(len(v) - 1)]
		URLsFound = append(URLsFound, currentURL, currentURL)
	}
	return URLsFound
}
