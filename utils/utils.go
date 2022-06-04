package utils

import (
	"crypto/sha256"
)

type HashValue []byte

func HashSum(a []HashValue) (sum HashValue) {
	h := sha256.New()
	for _, v := range a {
		h.Write(v)
	}
	h.Sum(sum)
	return
}
