package utils

import (
	"crypto/sha256"
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"
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

func CopyFile(src, dst string) (int64, error) {
	// log.Printf("CopyFile %s %s", src, dst)
	if src == dst {
		return 0, fmt.Errorf("same path")
	}
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func ReaderChecksum(reader io.Reader) Checksum {
	hash := sha256.New()
	if _, err := io.Copy(hash, reader); err != nil {
		return Checksum{}
	}
	var b = hash.Sum(nil)
	if len(b) != 32 {
		panic(b)
	}
	return *(*Checksum)(b)
}

// SHA256 hash for file content.
// for any error, return empty hash
func FileChecksum(name string) Checksum {
	f, err := os.Open(name)
	if err != nil {
		return Checksum{}
	}
	defer f.Close()

	return ReaderChecksum(f)
}

// comparable
type Checksum [32]byte

func (r Checksum) String() string {
	s := ""
	for _, v := range r {
		s += fmt.Sprintf("%02x", v)
	}
	return s
}

type LangTag int

const (
	Lcpp LangTag = iota
	Lcpp11
	Lcpp14
	Lcpp17
	Lcpp20
	Lpython2
	Lpython3
	Lgo
	Ljava
	Lc
)

type CtntType int

const (
	Cplain CtntType = iota
	Cbinary
	Csource
)

func init() {
	rand.Seed(time.Now().Unix())
}
