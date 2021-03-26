package crawl

import (
	"net/url"
	"os"
	"path/filepath"

	"github.com/OrlovM/go-workerpool"
	"github.com/urfave/cli/v2"
)

type urlSlice []*url.URL
type tasksSlice []workerpool.Task

//Checks if urlSlice contains Page.URL
func (s urlSlice) Contains(p *Page) bool {
	for _, u := range s {
		if *p.URL == *u {
			return true
		}
	}
	return false
}

//Trims the first task in tasksSlice
func (s *tasksSlice) TrimFirst() {
	*s = (*s)[1:]
}

func CreateFile(ctx *cli.Context) (*os.File, error) {
	dir, file := filepath.Split(ctx.String("filepath"))
	if err := os.MkdirAll(filepath.Dir(dir), 0770); err != nil {
		return nil, err
	}
	if file == "" {
		file = ctx.App.Compiled.Format("2006-01-02T15:04:05")
	}
	fp := filepath.Join(dir, file)
	return os.Create(fp)

}
