package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/sshwy/yaoj-core/pkg/migrator"
)

var isUoj bool
var srcDir string
var destDir string

func main() {
	flag.Parse()

	err := os.MkdirAll(destDir, os.ModePerm)
	if err != nil {
		fmt.Printf("error: %s", err.Error())
		return
	}
	fmt.Printf("output to %q\n", destDir)

	var mig migrator.Migrator
	if isUoj {
		mig = migrator.Uoj{}
	} else {
		fmt.Printf("type not specified\n")
		return
	}

	_, err = mig.Migrate(srcDir, destDir)
	if err != nil {
		panic(err)
	}
	fmt.Printf("done.")
}

func init() {
	flag.StringVar(&srcDir, "src", "", "source directory")
	flag.StringVar(&destDir, "output", "", "output directory")
	flag.BoolVar(&isUoj, "uoj", false, "migrate from uoj problem data")
}
