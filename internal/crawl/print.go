package crawl

import (
	"fmt"
)

type printer struct {
	Print bool
}

func newPrinter(print bool) *printer {
	return &printer{Print: print}
}

func (p *printer) InBase(pg *Page) {
	if p.Print {
		fmt.Println("URL:", pg.URL, "Source:", pg.Source, " is already in base")
	}
}

func (p *printer) Found(pg *Page) {
	if p.Print {
		fmt.Println("Found new URL:", pg.URL.String(), "Source:", pg.Source, "Depth:", pg.Depth, "Status code:", pg.StatusCode)
	}
}

func (p *printer) Error(errs []error) {
	if p.Print {
		for err := range errs {
			fmt.Println(err)
		}
	}
}
