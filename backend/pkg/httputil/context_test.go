package httputil

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func ExampleRequestWithContextSet() {
	req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080?age=1&name=john", nil)

	// passing data into request context
	req = RequestWithContextSet(req, "k1", "v1")
	// get value from request
	v := GetHTTPContext(req.Context(), "k1").(string)
	fmt.Println(v)
	// Output: v1
}

func TestRequestWithContextSet(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080?age=1&name=john", nil)

	// passing data into request context
	req = RequestWithContextSet(req, "k1", "v1")
	req = RequestWithContextSet(req, "k2", "v2")

	v := GetHTTPContext(req.Context(), "k1").(string)
	assert.Equal(t, "v1", v)
	v = GetHTTPContext(req.Context(), "k2").(string)
	assert.Equal(t, "v2", v)

	req = RequestWithContextSet(req, "k2", "v22", "k3", "v3")
	v = GetHTTPContext(req.Context(), "k2").(string)
	assert.Equal(t, "v22", v)
	v = GetHTTPContext(req.Context(), "k3").(string)
	assert.Equal(t, "v3", v)
}
