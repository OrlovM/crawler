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

func (p *printer) InBase(pg *crawler.Page) {
	if p.Print == true {
		fmt.Println("URL:", pg.URL, "Source:", pg.Source, " is already in base")
	}
}

func (p *printer) Found(pg *crawler.Page) {
	if p.Print == true {
		fmt.Println("Found new URL:", pg.URL.String(), "Source:", pg.Source, "Depth:", pg.Depth, "Status code:", pg.StatusCode)
	}
}

func (p *printer) Error(err error) {
	if p.Print == true {
		fmt.Println(err)
	}
}
