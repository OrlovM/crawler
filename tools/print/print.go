package print

import (
	"crawler/internal/crawler"
	"fmt"
)

type printer struct {
	Print bool
}

func NewPrinter(print bool) *printer {
	return &printer{Print: print}
}

func (p *printer) InBase(currentPage *crawler.Page) {
	if p.Print == true {
		fmt.Println("URL:", currentPage.URL, "Source:", currentPage.Source, " is already in base")
	}
}

func (p *printer) Found(currentPage *crawler.Page) {
	if p.Print == true {
		fmt.Println("Found new URL:", currentPage.URL, "Source:", currentPage.Source, "Depth:", currentPage.Depth, "Status code:", currentPage.StatusCode)
	}
}
