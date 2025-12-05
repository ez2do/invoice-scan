package httputil

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"invoice-scan/backend/pkg/errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

const errText = "error response_test"

func TestRespondJSON(t *testing.T) {
	type args struct {
		payload interface{}
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "string",
			args: args{payload: "gmicro"},
		},
		{
			name: "map",
			args: args{payload: map[string]string{"key": "value"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			RespondJSON(rr, http.StatusOK, tt.args.payload)
			rrb := rr.Body.Bytes()

			var body interface{}
			err := json.Unmarshal(rrb, &body)
			assert.Nil(t, err)
			assert.Equal(t, "application/json; charset=utf-8", rr.Header().Get("Content-Type"))
		})
	}
}

func TestResponseWrapJSON(t *testing.T) {
	type args struct {
		payload interface{}
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "string",
			args: args{payload: "gmicro"},
		},
		{
			name: "map",
			args: args{payload: map[string]string{"key": "value"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			ResponseWrapJSON(rr, http.StatusOK, tt.args.payload, true, "")
			rrb := rr.Body.Bytes()

			var body interface{}
			err := json.Unmarshal(rrb, &body)
			assert.Nil(t, err)
			assert.Equal(t, "application/json; charset=utf-8", rr.Header().Get("Content-Type"))
		})
	}
}

func TestRespondError(t *testing.T) {
	var body interface{}

	rr := httptest.NewRecorder()
	RespondError(rr, http.StatusOK, errText)
	rrb := rr.Body.Bytes()

	err := json.Unmarshal(rrb, &body)
	assert.Nil(t, err)

	want := map[string]interface{}{"error": errText}
	assert.Equal(t, body, want)
}

func TestResponseWrapError(t *testing.T) {
	var body interface{}

	rr := httptest.NewRecorder()
	ResponseWrapError(rr, http.StatusOK, errText)
	rrb := rr.Body.Bytes()

	err := json.Unmarshal(rrb, &body)
	assert.Nil(t, err)

	want := map[string]interface{}{"message": errText, "success": false}
	assert.Equal(t, body, want)
}

func TestRespondIError(t *testing.T) {
	err := errors.NewHTTPErrorCode(http.StatusBadRequest, "error_code", errText)

	rr := httptest.NewRecorder()
	RespondIError(rr, err)
	rrb := rr.Body.String()
	wantB := "{\"error_code\":\"error_code\",\"error_id\":\"" + err.ID() + "\",\"error_summary\":\"" + errText + "\"}"
	assert.Equal(t, wantB, rrb)
	assert.Equal(t, "error_code", rr.Header().Get("X-Error-Code"))
}

func TestResponseWrapIError(t *testing.T) {
	err := errors.NewHTTPErrorCode(http.StatusBadRequest, "error_code", errText)

	rr := httptest.NewRecorder()
	ResponseWrapIError(rr, err)
	rrb := rr.Body.String()
	wantB := "{\"success\":false,\"message\":\"error response_test\"," +
		"\"detail\":{\"error_code\":\"error_code\"," +
		"\"error_id\":\"" + err.ID() + "\",\"error_summary\":\"error response_test\"}}"
	assert.Equal(t, wantB, rrb)
	assert.Equal(t, "error_code", rr.Header().Get("X-Error-Code"))
}

func TestRespondMessage(t *testing.T) {
	rr := httptest.NewRecorder()

	RespondMessage(rr, http.StatusOK, errText)
	rrb := rr.Body.String()

	want := "{\"message\":\"" + errText + "\"}"
	assert.Equal(t, rrb, want)
}

func TestRespondString(t *testing.T) {
	rr := httptest.NewRecorder()
	textRS := []byte("Error RespondString")
	RespondString(rr, http.StatusOK, "Error RespondString")
	rrb := rr.Body.Bytes()

	assert.Equal(t, rrb, textRS)
}

func TestResponseRedirect(t *testing.T) {
	rr := httptest.NewRecorder()
	ResponseRedirect(rr, http.StatusMovedPermanently, "http://localhost:8080/redirect/to")
	header := rr.Header()
	loc := header.Get("location")
	assert.Equal(t, loc, "http://localhost:8080/redirect/to")
}
