package tools

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

func CreateFile(ctx *cli.Context) (*os.File, error) {
	dir, file := filepath.Split(ctx.String("filepath"))
	if err := os.MkdirAll(filepath.Dir(dir), 0770); err != nil {
		return nil, err
	}
	if file == "" {
		file = fmt.Sprintf("%v%v", ctx.App.Compiled.Format("2006-01-02T15:04:05"), ".json")
	}
	fp := filepath.Join(dir, file)
	return os.Create(fp)

}
