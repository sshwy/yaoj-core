package problem_test

import (
	"testing"

	"github.com/k0kubun/pp/v3"
	"github.com/sshwy/yaoj-core/problem"
)

func TestProbDtgp(t *testing.T) {
	group, err := problem.LoadDtgp("testdata/prob/datagroup/testcase")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("name: ", group.Name())
	t.Log(pp.Sprint(group.Records()))
	err = group.AddField("config")
	if err != nil {
		t.Error(err)
		return
	}
	// t.Log(pp.Sprint(group.Records()))
	err = group.NewRecord()
	if err != nil {
		t.Error(err)
		return
	}
	// t.Log(pp.Sprint(group.Records()))
	err = group.RemoveField("config")
	if err != nil {
		t.Error(err)
		return
	}
	// t.Log(pp.Sprint(group.Records()))
	err = group.RemoveRecord(2)
	if err != nil {
		t.Error(err)
		return
	}
	// t.Log(pp.Sprint(group.Records()))
	err = group.AlterValue(1, "input", "testdata/data.in")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestNew(t *testing.T) {
	dir := t.TempDir()
	prob, err := problem.New(dir)
	if err != nil {
		t.Error(err)
		return
	}
	prob.SetStmt([]byte("test"))
}

func TestLoad(t *testing.T) {
	prob, err := problem.Load("testdata/prob")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(pp.Sprint(prob))
}
