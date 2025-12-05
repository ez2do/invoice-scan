package httputil

import (
	"context"
	"invoice-scan/backend/pkg/log"
	"net/http"

	"github.com/spf13/cast"
)

type httpContextKey int

const httpContextID httpContextKey = iota

// RequestWithContextSet return request with http context that contains new k-v data,
// It's not thread-safe for using on same r in multiple go routines (cause concurrent writing panic)
func RequestWithContextSet(r *http.Request, kvs ...interface{}) *http.Request {
	rCtx := r.Context()
	// already exist
	if vmap, ok := rCtx.Value(httpContextID).(map[string]interface{}); ok {
		for i := 0; i < len(kvs); {
			if i == len(kvs)-1 {
				log.Warnf("ignore key without value: %v", kvs[i])
				break
			}
			vmap[cast.ToString(kvs[i])] = kvs[i+1]
			i += 2
		}
		return r
	}

	// initialize new instance
	m := make(map[string]interface{})
	rCtx = context.WithValue(rCtx, httpContextID, m)
	for i := 0; i < len(kvs); {
		if i == len(kvs)-1 {
			log.Warnf("ignore key without value: %v", kvs[i])
			break
		}
		m[cast.ToString(kvs[i])] = kvs[i+1]
		i += 2
	}
	return r.WithContext(rCtx)
}

// GetHTTPContext return value of k inside request http context
func GetHTTPContext(ctx context.Context, k string) interface{} {
	if vmap, ok := ctx.Value(httpContextID).(map[string]interface{}); ok {
		return vmap[k]
	}

	return nil
}
