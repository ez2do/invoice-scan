package pkg

import (
	"crypto/rand"
	"invoice-scan/backend/pkg/log"
	"math/big"
)

var (
	randLetters = []rune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
)

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(n int) []byte {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		log.Fatalf("Error when gen rand bytes: %v", err)
		return nil
	}

	return b
}

// GenerateRandomKey generate random letters with fixed length
func GenerateRandomKey(n int) string {
	return GenerateRandomWithChars(n, randLetters)
}

// GenerateRandomNum generate random number in range [0,n)
func GenerateRandomNum(n int64) uint64 {
	nBig, err := rand.Int(rand.Reader, big.NewInt(n))
	if err != nil {
		log.Fatalf("Fail to generate random number in range 0-%v: %v", n, err)
	}
	return nBig.Uint64()
}

// GenerateRandomWithChars generate random strings with fixed length from charset
func GenerateRandomWithChars(n int, chars []rune) string {
	var esp byte
	resp, err := rand.Int(rand.Reader, big.NewInt(256))
	if err != nil {
		log.Fatalf("Error when rand int: %v", err)
	}
	esp = byte(resp.Int64() % int64(len(chars)))

	// add amount of resp in range [0, 256 mod len(chars)] to complement for random bytes
	rbs := GenerateRandomBytes(n)
	b := make([]rune, n)
	for i, rb := range rbs {
		b[i] = chars[(esp+rb)%byte(len(chars))]
	}
	return string(b)
}
