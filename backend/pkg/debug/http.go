package debug

import (
	"fmt"
	"invoice-scan/backend/pkg"
	"invoice-scan/backend/pkg/log"
	"io"
	"net/http"
	"time"
)

type httpRoundTripper struct {
	origin http.RoundTripper
}

func NewHTTPRoundTripper(origin http.RoundTripper) http.RoundTripper {
	return &httpRoundTripper{
		origin: origin,
	}
}

func (r *httpRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// Get request body
	var requestBody io.ReadCloser
	if req.Body != nil && req.GetBody != nil {
		requestBody, _ = req.GetBody()
		if requestBody != nil {
			defer requestBody.Close()
		}
	}

	startTime := time.Now()
	response, err := r.origin.RoundTrip(req)

	var requestBodyBytes []byte
	if requestBody != nil {
		requestBodyBytes, _ = io.ReadAll(requestBody)
	}

	var statusCode int
	if response != nil {
		statusCode = response.StatusCode
	}

	processTime := time.Since(startTime)
	log.Infow(
		fmt.Sprintf("HTTP request completed (took: %s)", processTime),
		"error", err,
		"context", pkg.ToJSONString(map[string]interface{}{
			"processTime": processTime.Seconds(),
			"status_code": statusCode,
			"request_url": req.URL.String(),
			"headers":     req.Header,
			"method":      req.Method,
			"body":        string(requestBodyBytes),
		}),
	)

	return response, err
}
