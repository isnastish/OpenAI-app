package hashutils

import (
	"crypto/sha256"
	"fmt"
)

func NewSha256(data []byte) string {
	hash := sha256.New()
	hash.Write(data)
	res := hash.Sum(nil)
	return fmt.Sprintf("%x", res)
}
