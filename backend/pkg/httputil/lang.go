package httputil

import (
	"golang.org/x/text/language"
	"net/http"
	"strings"
)

// GetAcceptLang get list of accept language from request in order
func GetAcceptLang(r *http.Request) []language.Tag {
	acceptLangs := r.Header.Values(HeaderAcceptLanguage)
	var langs []language.Tag
	for _, al := range acceptLangs {
		l := strings.Split(al, ";")[0]
		if ml, err := language.Parse(l); err == nil {
			langs = append(langs, ml)
		}
	}
	return langs
}
