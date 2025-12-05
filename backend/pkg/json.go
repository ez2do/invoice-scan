package pkg

import (
	"encoding/json"
	"invoice-scan/backend/pkg/log"
)

func ToJSONString(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		log.Errorw("Cannot marshal json", "error", err)
		return ""
	}

	return string(b)
}
