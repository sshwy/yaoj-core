package problem_test

import (
	"path"
	"testing"

	"github.com/bitfield/script"
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
	rcd, err := group.NewRecord()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(pp.Sprint(rcd))
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
	err = group.Record(1).AlterValue("input", "testdata/data.in")
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
	prob.SetStmt([]byte("a plus b problem!"))
	testcase, err := prob.NewDataGroup("testcase")
	if err != nil {
		t.Error(err)
		return
	}
	script.Echo("1 2").WriteFile(path.Join(dir, "a.in"))
	script.Echo("3").WriteFile(path.Join(dir, "a.ans"))

	testcase.AddField("input")
	testcase.AddField("answer")
	r0, _ := testcase.NewRecord()
	r0.AlterValue("input", path.Join(dir, "a.in"))
	r0.AlterValue("answer", path.Join(dir, "a.ans"))
	r1, _ := testcase.NewRecord()
	r1.AlterValue("input", path.Join(dir, "a.in"))
	r1.AlterValue("answer", path.Join(dir, "a.ans"))

	script.Echo(`
#include<iostream>
using namespace std;
int main () { int a, b; cin >> a >> b; cout << a + b << endl; return 0; }
	`).WriteFile(path.Join(dir, "src.cpp"))

	_submission, err := prob.NewDataGroup("submission")
	if err != nil {
		t.Error(err)
		return
	}
	_submission.AddField("source")
	_submission.NewRecord()

	script.Echo("1000 1000 104857600 104857600 104857600 104857600 10").WriteFile(path.Join(dir, "cpl.txt"))
	script.Echo("#!/bin/env bash\ng++ $1 -o $2 -O2").WriteFile(path.Join(dir, "script.sh"))

	option, err := prob.NewDataGroup("option")
	if err != nil {
		t.Error(err)
		return
	}
	option.AddField("limitation")
	option.AddField("compilescript")
	o0, _ := option.NewRecord()
	o0.AlterValue("limitation", path.Join(dir, "cpl.txt"))
	o0.AlterValue("compilescript", path.Join(dir, "script.sh"))

	err = prob.SetWkflGraph([]byte(`{"Node":[{"ProcName":"compiler","InEdge":[{"Bound":null,"Label":"source"},{"Bound":null,"Label":"script"}],"OutEdge":[{"Bound":{"Index":1,"LabelIndex":0},"Label":"result"},{"Bound":null,"Label":"log"},{"Bound":null,"Label":"judgerlog"}]},{"ProcName":"runner:stdio","InEdge":[{"Bound":{"Index":0,"LabelIndex":0},"Label":"executable"},{"Bound":null,"Label":"stdin"},{"Bound":null,"Label":"limit"}],"OutEdge":[{"Bound":{"Index":2,"LabelIndex":0},"Label":"stdout"},{"Bound":null,"Label":"stderr"},{"Bound":null,"Label":"judgerlog"}]},{"ProcName":"checker:hcmp","InEdge":[{"Bound":{"Index":1,"LabelIndex":0},"Label":"out"},{"Bound":null,"Label":"ans"}],"OutEdge":[{"Bound":null,"Label":"result"}]}],"Inbound":{"option":{"compilescript":{"Index":0,"LabelIndex":1},"limitation":{"Index":1,"LabelIndex":2}},"submission":{"source":{"Index":0,"LabelIndex":0}},"testcase":{"answer":{"Index":2,"LabelIndex":1},"input":{"Index":1,"LabelIndex":1}}}}`))
	if err != nil {
		t.Error(err)
		return
	}
	pp.Print(prob)

	submdir := t.TempDir()
	submission, err := problem.LoadDtgp(submdir)
	if err != nil {
		t.Error(err)
		return
	}
	for _, field := range _submission.Fields() { // aka "source"
		submission.AddField(field)
	}
	s0, _ := submission.NewRecord()
	s0.AlterValue("source", path.Join(dir, "src.cpp"))

	testcases := prob.Testcase()
	for i, testcase := range testcases {
		t.Log(i, " ", testcase.ID())
		t.Log(pp.Sprint(testcase.Run(t.TempDir(), *submission.Record(0).PathMap(), 100)))
	}
}

func TestLoad(t *testing.T) {
	_, err := problem.Load("testdata/prob")
	if err != nil {
		t.Error(err)
		return
	}
	// t.Log(pp.Sprint(prob))
}
