package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/sshwy/yaoj-core/pkg/migrator"
)

var isUoj bool
var srcDir string
var destDir string
var dumpFile string

func Main() error {
	err := os.MkdirAll(destDir, os.ModePerm)
	if err != nil {
		return err
	}
	fmt.Printf("output to %q\n", destDir)

	var mig migrator.Migrator
	if isUoj {
		mig = migrator.Uoj{}
	} else {
		return fmt.Errorf("type not specified")
	}

	if dumpFile == "" {
		_, err = mig.Migrate(srcDir, destDir)
		if err != nil {
			return err
		}
	} else {
		dir, err := os.MkdirTemp(os.TempDir(), "yaoj-migrator-******")
		if err != nil {
			return err
		}
		prob, err := mig.Migrate(srcDir, dir)
		if err != nil {
			return err
		}
		err = prob.DumpFile(path.Join(destDir, dumpFile))
		if err != nil {
			return err
		}
		err = os.RemoveAll(dir)
		if err != nil {
			return err
		}
	}
	fmt.Printf("done.")
	return nil
}

func main() {
	flag.Parse()

	err := Main()
	if err != nil {
		fmt.Printf("[error]: %s\n", err.Error())
		flag.Usage()
		return
	}
}

func init() {
	flag.StringVar(&srcDir, "src", "", "source directory")
	flag.StringVar(&destDir, "output", "", "output directory")
	flag.StringVar(&dumpFile, "dump", "", "output a zip archive with given name")
	flag.BoolVar(&isUoj, "uoj", false, "migrate from uoj problem data")
}
