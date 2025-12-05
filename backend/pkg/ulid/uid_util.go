package ulid

import (
	crand "crypto/rand"
	"github.com/oklog/ulid/v2"
)

// GenerateULID generate unique ID with ULID format
func GenerateULID() string {
	return ulid.MustNew(ulid.Now(), crand.Reader).String()
}
