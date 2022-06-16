package problem_test

import (
	"log"

	"github.com/sshwy/yaoj-core/problem"
)

func ExampleProbDtgp() {
	dir := "/path/to/an/empty/dir"
	// Create a datagroup by load an empty dir
	group, err := problem.LoadDtgp(dir)
	if err != nil {
		log.Fatal(err)
	}
	group.AddField("username")
	group.AddField("password")
	record, err := group.NewRecord()
	if err != nil {
		log.Fatal(err)
	}
	record.AlterValue("username", "path/to/a/text/file")
	record.AlterValue("password", "path/to/another/text/file")
	log.Print(record.Map)
	log.Print(record.PathMap())
}
