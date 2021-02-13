package crawl

import (
	"net/url"
	"testing"
)

type testTask struct {
	i int
}

func (t *testTask) Process() {
	t.i = 0
}

func TestTasksSlice_TrimFirst(t *testing.T) {
	slice := TasksSlice{&testTask{0}, &testTask{1}, &testTask{2}, &testTask{3}}
	slice.TrimFirst()
	want := TasksSlice{&testTask{1}, &testTask{2}, &testTask{3}}
	for i, ts := range slice{
		if ts.(*testTask).i != want[i].(*testTask).i {
			t.Errorf("Slice %d is not equal to want %d", i, i)
		}
	}

}

func TestURLSlice_Contains(t *testing.T) {
	first, _ := url.Parse("http://ya.ru")
	second, _ := url.Parse("https://www.google.ru/")
	third, _ := url.Parse("https://mail.google.com/")
	fourth, _ := url.Parse("https://golang.org/")
	fifth, _ := url.Parse("https://www.youtube.com/")
	type args struct {
		p *Page
	}
	tests := []struct {
		name string
		s    URLSlice
		args args
		want bool
	}{
		{name: "contains", s: URLSlice{first, second,third, fourth, fifth}, args: args{&Page{URL: fifth}}, want: true},
		{name: "notContains", s: URLSlice{first, second,third, fourth}, args: args{&Page{URL: fifth}}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Contains(tt.args.p); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}