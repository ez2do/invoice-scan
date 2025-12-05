package httputil

import (
	"errors"
	"fmt"
	"invoice-scan/backend/pkg/log"
	"net/http"
	"strings"
)

type ResponseOption func(w http.ResponseWriter)

func WithHeaders(kvs ...string) ResponseOption {
	return func(w http.ResponseWriter) {
		for i := 0; i < len(kvs); {
			if i == len(kvs)-1 {
				log.Errorw("ignore key without values", "error", errors.New("key without value"),
					"context", fmt.Sprintf("key: %s", kvs[i]))
				break
			}
			key, val := kvs[i], kvs[i+1]
			w.Header().Set(key, val)
			i += 2
		}
	}
}

func WithContentLanguage(lang string) ResponseOption {
	return func(w http.ResponseWriter) {
		langs := strings.Split(lang, ",")
		for _, l := range langs {
			w.Header().Add(HeaderContentLanguage, l)
		}
	}
}
