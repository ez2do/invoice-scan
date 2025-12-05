package ulid

import (
	"github.com/oklog/ulid/v2"
	"log"
	"testing"
)

func TestGenerateULID(t *testing.T) {
	uniqueID := GenerateULID()
	log.Print(uniqueID)

	ULIDObj := ulid.MustParse(uniqueID)
	log.Printf("%v %v", ULIDObj.Time(), ULIDObj.Entropy())
}
