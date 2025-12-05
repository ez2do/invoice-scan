package httputil

import (
	"golang.org/x/text/language"
	"net/http"
	"strings"
	"testing"
)

func buildGetRequest(url string, kv http.Header) *http.Request {
	r, _ := http.NewRequest(http.MethodGet, url, nil)
	r.Header = kv.Clone()
	return r
}

func TestGetAcceptLang(t *testing.T) {
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "request with accept language",
			args: args{r: buildGetRequest("localhost", http.Header{
				HeaderAcceptLanguage: []string{"en-US"},
			})},
			want: "en-US",
		},
		{
			name: "request with accept language without variant",
			args: args{r: buildGetRequest("localhost", http.Header{
				HeaderAcceptLanguage: []string{"vi"},
			})},
			want: "vi",
		},
		{
			name: "request without language",
			args: args{r: buildGetRequest("https://google.com", http.Header{})},
			want: "",
		},
		{
			name: "request with wrong language",
			args: args{r: buildGetRequest("localhost", http.Header{
				HeaderAcceptLanguage: []string{"xon-exist"},
			})},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetAcceptLang(tt.args.r); langStr(got) != tt.want {
				t.Errorf("GetAcceptLang() = %v, want %v", got, tt.want)
			}
		})
	}
}

func langStr(got []language.Tag) string {
	var langs []string
	for _, l := range got {
		langs = append(langs, l.String())
	}
	return strings.Join(langs, ",")
}
