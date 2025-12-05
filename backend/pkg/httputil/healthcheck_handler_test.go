package httputil

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"invoice-scan/backend/pkg"

	"net/http"
	"net/http/httptest"
	"testing"
)

type FakeServiceHealth struct {
	status int
}

func NewFakeServiceHealth(status int) *FakeServiceHealth {
	return &FakeServiceHealth{status: status}
}

func (r *FakeServiceHealth) Check() error {
	if r.status == http.StatusOK {
		return nil
	}
	return errors.New("Error check health service")
}

func (r *FakeServiceHealth) Name() string {
	return "FakeCheck"
}

func prepareHealthMock(status int) http.Handler {
	fsh := NewFakeServiceHealth(status)
	hc := NewHealthCheck(fsh)

	return hc
}

func TestHealthCheckHealthyHandler_ServeHTTP(t *testing.T) {
	rr := httptest.NewRecorder()

	hc := prepareHealthMock(http.StatusOK)

	req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080?age=1&name=john", nil)

	hc.ServeHTTP(rr, req)
	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, rr.Body.String(), "{\"message\":\"Service is healthy\",\"status\":\"healthy\"}")
}

func TestHealthCheckUnhealthyHandler_ServeHTTP(t *testing.T) {
	rr := httptest.NewRecorder()
	hc := prepareHealthMock(http.StatusServiceUnavailable)

	req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080?age=1&name=john", nil)
	hc.ServeHTTP(rr, req)

	rrb := rr.Body.Bytes()
	var body map[string]interface{}
	err := json.Unmarshal(rrb, &body)
	assert.Nil(t, err)

	assert.Equal(t, rr.Code, http.StatusServiceUnavailable)
	assert.Equal(t, body["status"], string(pkg.ServiceUnhealthy))
	assert.NotNil(t, body["errors"])
}
