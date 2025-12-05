package httputil

import (
	"encoding/json"
	"github.com/spf13/cast"
	"invoice-scan/backend/pkg/errors"
	"invoice-scan/backend/pkg/locale"
	"net/http"
	"strconv"
)

const (
	XErrorIDHeader   = "X-Error-ID"
	XErrorCodeHeader = "X-Error-Code"
)

type baseResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Detail  interface{} `json:"detail,omitempty"`
}

func responseRawJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.WriteHeader(statusCode)
	_, _ = w.Write(data)
}

// RespondJSON -- makes the tracking_resp with payload as json format
func RespondJSON(w http.ResponseWriter, httpStatusCode int, payload interface{}, opts ...ResponseOption) {
	for _, opt := range opts {
		opt(w)
	}
	responseRawJSON(w, httpStatusCode, payload)
}

// RespondError -- makes the error tracking_resp with payload as json format
func RespondError(w http.ResponseWriter, httpStatusCode int, message string, opts ...ResponseOption) {
	for _, opt := range opts {
		opt(w)
	}
	message = locale.TL(w.Header().Get(HeaderContentLanguage), message)
	responseRawJSON(w, httpStatusCode, map[string]string{"error": message})
}

// RespondIError -- makes the error tracking_resp with payload as json format via IError instance
func RespondIError(w http.ResponseWriter, ierr *errors.IError, opts ...ResponseOption) {
	opts = append(opts, WithHeaders(XErrorIDHeader, ierr.ID()), WithHeaders(XErrorCodeHeader, ierr.Code()))
	for _, opt := range opts {
		opt(w)
	}
	errDetail := translateIError(ierr, w.Header().Get(HeaderContentLanguage))
	responseRawJSON(w, ierr.HTTPCode(), errDetail)
}

// RespondMessage -- makes the message tracking_resp with payload as json format
func RespondMessage(w http.ResponseWriter, httpStatusCode int, message string, opts ...ResponseOption) {
	for _, opt := range opts {
		opt(w)
	}
	message = locale.TL(w.Header().Get(HeaderContentLanguage), message)
	responseRawJSON(w, httpStatusCode, map[string]string{"message": message})
}

// RespondString -- makes the string
func RespondString(w http.ResponseWriter, httpStatusCode int, message string, opts ...ResponseOption) {
	for _, opt := range opts {
		opt(w)
	}
	message = locale.TL(w.Header().Get(HeaderContentLanguage), message)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(len([]byte(message))))
	w.WriteHeader(httpStatusCode)
	_, _ = w.Write([]byte(message))
}

// RespondHTML -- response HTML content
func RespondHTML(w http.ResponseWriter, httpStatusCode int, content string, opts ...ResponseOption) {
	for _, opt := range opts {
		opt(w)
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(len([]byte(content))))
	w.WriteHeader(httpStatusCode)
	_, _ = w.Write([]byte(content))
}

// ResponseRedirect -- makes redirect response
func ResponseRedirect(w http.ResponseWriter, httpStatusCode int, location string, opts ...ResponseOption) {
	for _, opt := range opts {
		opt(w)
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Location", location)
	w.WriteHeader(httpStatusCode)
}

// RespondWrapMessage -- makes the message tracking_resp with payload as json format
// then wrap it into GHTK response format
func RespondWrapMessage(w http.ResponseWriter, httpStatusCode int, message string, opts ...ResponseOption) {
	for _, opt := range opts {
		opt(w)
	}
	base := &baseResponse{
		Success: true,
		Message: locale.TL(w.Header().Get(HeaderContentLanguage), message),
		Data:    nil,
		Detail:  nil,
	}
	responseRawJSON(w, httpStatusCode, base)
}

// ResponseWrapJSON -- makes the tracking_resp with payload as json format then wrap it into GHTK response format
func ResponseWrapJSON(w http.ResponseWriter, httpStatusCode int, payload interface{},
	success bool, message string, opts ...ResponseOption) {
	for _, opt := range opts {
		opt(w)
	}
	base := &baseResponse{
		Success: success,
		Message: locale.TL(w.Header().Get(HeaderContentLanguage), message),
		Data:    payload,
		Detail:  nil,
	}
	responseRawJSON(w, httpStatusCode, base)
}

// ResponseWrapJSONError -- makes the tracking_resp with detail as json format then wrap it into GHTK response format
func ResponseWrapJSONError(w http.ResponseWriter, httpStatusCode int, detail interface{},
	message string, opts ...ResponseOption) {
	for _, opt := range opts {
		opt(w)
	}
	base := &baseResponse{
		Success: false,
		Message: locale.TL(w.Header().Get(HeaderContentLanguage), message),
		Data:    nil,
		Detail:  detail,
	}
	responseRawJSON(w, httpStatusCode, base)
}

// ResponseWrapSuccessJSON -- makes the tracking_resp with payload as json format then wrap it into GHTK response format
func ResponseWrapSuccessJSON(w http.ResponseWriter, httpStatusCode int, payload interface{}, opts ...ResponseOption) {
	ResponseWrapJSON(w, httpStatusCode, payload, true, "", opts...)
}

// ResponseWrapFailJSON -- makes the tracking_resp with payload as json format then wrap it into GHTK response format
func ResponseWrapFailJSON(
	w http.ResponseWriter, httpStatusCode int, payload interface{}, mess string, opts ...ResponseOption) {
	ResponseWrapJSON(w, httpStatusCode, payload, false, mess, opts...)
}

// ResponseWrapError -- makes the error tracking_resp with payload as json format then wrap it into GHTK response format
func ResponseWrapError(w http.ResponseWriter, httpStatusCode int, message string, opts ...ResponseOption) {
	for _, opt := range opts {
		opt(w)
	}
	base := &baseResponse{
		Success: false,
		Message: locale.TL(w.Header().Get(HeaderContentLanguage), message),
		Data:    nil,
		Detail:  nil,
	}
	responseRawJSON(w, httpStatusCode, base)
}

// ResponseWrapIError -- makes the error tracking_resp with payload as json format via IError instance
// then wrap it into GHTK response format
func ResponseWrapIError(w http.ResponseWriter, ierr *errors.IError, opts ...ResponseOption) {
	// check if content-language is set
	opts = append(opts, WithHeaders(XErrorIDHeader, ierr.ID()), WithHeaders(XErrorCodeHeader, ierr.Code()))
	for _, opt := range opts {
		opt(w)
	}
	errDetail := translateIError(ierr, w.Header().Get(HeaderContentLanguage))
	base := &baseResponse{
		Success: false,
		Message: cast.ToString(errDetail["error_summary"]),
		Data:    nil,
		Detail:  errDetail,
	}
	responseRawJSON(w, ierr.HTTPCode(), base)
}

func translateIError(ierr *errors.IError, lang string) map[string]interface{} {
	errDetail := ierr.JSON()
	msg := cast.ToString(errDetail["error_summary"])
	errDetail["error_summary"] = locale.TL(lang, msg, locale.WithTplData(ierr.Data()))
	return errDetail
}
