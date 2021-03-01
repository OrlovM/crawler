package crawl

import (
	"reflect"
	"strings"
	"testing"
)

const (
	htmlFragment = `<html>
		  <head>
		  <title> TestFragment </title>
		  </head>
		  <body bgcolor="pink">
		  <br>
		  <center><h2>This is a test <a href="http://test.com">URL </a></h2>
		  <br>
		  <p>Content here</p>
		  <a href="http://test2.com/test">Test 2 </a>
		  </center>
		  </body>
		  </html>`
)

func TestGetLinks(t *testing.T) {
	tests := []struct {
		body string
		want []string
	}{
		{body: "This is just a test string", want: []string{}},
		{body: `<a href="http://www.yandex.ru">ya.ru</a>`, want: []string{"http://www.yandex.ru"}},
		{body: htmlFragment, want: []string{"http://test.com", "http://test2.com/test"}},
	}

	for _, tt := range tests {
		r := strings.NewReader(tt.body)
		got := getLinks(r)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("not equal %v %v %d %d", got, tt.want, len(got), len(tt.want))

		}
	}
}
