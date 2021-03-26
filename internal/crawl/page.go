package crawl

import (
	"encoding/json"
	"net/url"
)

//Page is a representation of web page
type Page struct {
	URL        *url.URL
	Depth      int
	Source     string // URL where this Page was found
	StatusCode int    //0 if Page was not requested or request failed
}

//PagesSlice is a slice of Pages. Made for ability to implement methods on slice.
type PagesSlice []Page

//Implementation of sort interface
func (s PagesSlice) Len() int { return len(s) }

func (s PagesSlice) Less(i, j int) bool { return s[i].Depth < s[j].Depth }

func (s PagesSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

//MarshalJSON is Marshaler interface method
func (p *Page) MarshalJSON() ([]byte, error) {
	type Alias Page
	type JsonURL struct {
		Scheme      string        `json:"scheme,omitempty"`
		Opaque      string        `json:"opaque,omitempty"`
		User        *url.Userinfo `json:"user,omitempty"`
		Host        string        `json:"host,omitempty"`
		Path        string        `json:"path,omitempty"`
		RawPath     string        `json:"raw_path,omitempty"`
		ForceQuery  bool          `json:"force_query,omitempty"`
		RawQuery    string        `json:"raw_query,omitempty"`
		Fragment    string        `json:"fragment,omitempty"`
		RawFragment string        `json:"raw_fragment,omitempty"`
	}
	u := JsonURL{
		Scheme:      p.URL.Scheme,
		Opaque:      p.URL.Opaque,
		User:        p.URL.User,
		Host:        p.URL.Host,
		Path:        p.URL.Path,
		RawPath:     p.URL.RawPath,
		ForceQuery:  p.URL.ForceQuery,
		RawQuery:    p.URL.RawQuery,
		Fragment:    p.URL.Fragment,
		RawFragment: p.URL.RawFragment,
	}

	return json.Marshal(struct {
		URL        JsonURL `json:"url"`
		Depth      int     `json:"depth"`
		Source     string  `json:"source"`
		StatusCode int     `json:"status_code"`
	}{
		URL:        u,
		Depth:      p.Depth,
		Source:     p.Source,
		StatusCode: p.StatusCode,
	})
}
