package main

import (
	"flag"
	"log"
	"os"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/sshwy/yaoj-core/pkg/problem"
)

type Storage map[string]problem.Problem

var storage = Storage{}

func (r Storage) Has(checksum string) bool {
	_, ok := r[checksum]
	return ok
}
func (r Storage) Set(checksum string, prob problem.Problem) {
	r[checksum] = prob
}
func (r Storage) Get(checksum string) problem.Problem {
	return r[checksum]
}

var address string

func main() {
	flag.Parse()

	var cachedir = path.Join(os.TempDir(), "yaoj-judger-server-cache")
	os.RemoveAll(cachedir)
	if err := os.MkdirAll(cachedir, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	r.POST("/judge", Judge)
	r.POST("/sync", Sync)

	err := r.Run(address) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	if err != nil {
		logger.Fatal(err)
	}
}

func init() {
	flag.StringVar(&address, "listen", "localhost:3000", "listening address")
}

var logger = log.New(os.Stderr, "[judgeserver] ", log.LstdFlags|log.Lshortfile|log.Lmsgprefix)
