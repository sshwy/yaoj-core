package utils

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
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

type ByteValue int64

func (r ByteValue) String() string {
	num := float64(r)
	if num < 1000 {
		return fmt.Sprint(int64(num), "B")
	} else if num < 1e6 {
		return fmt.Sprintf("%.1f%s", num/1e3, "KB")
	} else if num < 1e9 {
		return fmt.Sprintf("%.1f%s", num/1e6, "MB")
	} else {
		return fmt.Sprintf("%.1f%s", num/1e9, "GB")
	}
}

func Map[T any, M any](s []T, f func(T) M) []M {
	var a []M = make([]M, len(s))
	for i, v := range s {
		a[i] = f(v)
	}
	return a
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
