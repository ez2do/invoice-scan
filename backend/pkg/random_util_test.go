package pkg

import (
	"github.com/stretchr/testify/assert"
	"invoice-scan/backend/pkg/log"
	"testing"
)

func TestGenerateRandomKey(t *testing.T) {
	var m = make(map[string]bool)
	var countDup = 0
	for i := 0; i < 1e5; i++ {
		t := GenerateRandomKey(6)
		if m[t] {
			countDup++
		}
		m[t] = true
	}
	log.Info("dup: ", countDup, ", m: ", len(m))
	assert.Less(t, float64(countDup)/float64(len(m)), 0.00002)
}

func TestGenerateRandomNum(t *testing.T) {
	var m = make(map[uint64]bool)
	var mDup = make(map[uint64]bool)
	for i := 0; i < 1e4; i++ {
		t := GenerateRandomNum(1e6)
		if m[t] {
			mDup[t] = true
		}
		m[t] = true
	}
	log.Info("dup: ", len(mDup), ", m: ", len(m))
	assert.Less(t, float64(len(mDup))/float64(len(m)), 0.01)
}

func BenchmarkGenerateRandomKey(b *testing.B) {
	for n := 0; n < b.N; n++ {
		GenerateRandomKey(64)
	}
}

func BenchmarkGenerateRandomNum(b *testing.B) {
	for n := 0; n < b.N; n++ {
		GenerateRandomNum(1e13)
	}
}

func BenchmarkGenerateRandomBytes(b *testing.B) {
	for n := 0; n < b.N; n++ {
		GenerateRandomBytes(64)
	}
}
