package httputil

type APIResult string

const (
	Success   APIResult = "success"
	Failure   APIResult = "failure"
	Skipped   APIResult = "skipped"
	Allow     APIResult = "allow"
	Deny      APIResult = "deny"
	Challenge APIResult = "challenge"
	Unknown   APIResult = "unknown"
)
