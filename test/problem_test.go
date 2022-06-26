package test_test

import (
	"path"
	"testing"

	"github.com/bitfield/script"
	"github.com/k0kubun/pp/v3"
	"github.com/sshwy/yaoj-core/pkg/problem"
)

var probData *problem.ProbData
var probDataDir string

func MakeProbData(t *testing.T) {
	dir := t.TempDir()
	var err error
	probData, err = problem.NewProbData(dir)
	if err != nil {
		t.Error(err)
		return
	}

	script.Echo("1 2").WriteFile(path.Join(dir, "a.in"))
	script.Echo("-1093908432").WriteFile(path.Join(dir, "a.ans"))
	script.Echo(`
#include<iostream>
using namespace std;

int main () { 
  int a, b; 
  cin >> a >> b;
  for(int i = 0; i < 100000000; i++) a += b, b += a;
  cout << a + b << endl;
  return 0;
}
	`).WriteFile(path.Join(dir, "src.cpp"))
	script.Echo("1000 1000 204857600 204857600 204857600 204857600 10").WriteFile(path.Join(dir, "cpl.txt"))
	script.Echo("#!/bin/env bash\nclang++ $1 -o $2").WriteFile(path.Join(dir, "script.sh"))
	script.Echo("# A + B Problem").WriteFile(path.Join(dir, "tmp.md"))

	probData.Fullscore = 100
	probData.Tests.Fields().Add("input")
	probData.Tests.Fields().Add("answer")
	probData.Tests.Fields().Add("_subtaskid")
	probData.Tests.Fields().Add("_score")
	r0 := probData.Tests.Records().New()
	r0["input"], err = probData.AddFile("a.in", path.Join(dir, "a.in"))
	if err != nil {
		t.Error(err)
		return
	}
	r0["answer"], err = probData.AddFile("a.ans", path.Join(dir, "a.ans"))
	if err != nil {
		t.Error(err)
		return
	}
	r0["_score"] = "average"

	r1 := probData.Tests.Records().New()
	r1["input"] = r0["input"]
	r1["answer"] = r0["answer"]
	r1["_score"] = "average"

	r2 := probData.Tests.Records().New()
	r2["input"] = r0["input"]
	r2["answer"] = r0["answer"]
	r2["_score"] = "average"

	r3 := probData.Tests.Records().New()
	r3["input"] = r0["input"]
	r3["answer"] = r0["answer"]
	r3["_score"] = "average"

	probData.Static.Fields().Add("limitation")
	probData.Static.Fields().Add("compilescript")
	o0 := probData.Static.Records().New()
	o0["limitation"] = "cpl.txt"
	o0["compilescript"] = "script.sh"

	probData.SetStmt("zh", "tmp.md")

	// net adjuestment
	err = probData.SetWkflGraph(wkflGraph.Serialize())
	if err != nil {
		t.Error(err)
		return
	}
	probData.Submission.Fields().Add("source")
	// pp.Print(prob)

	if err := probData.Export(probDataDir); err != nil {
		t.Error(err)
		return
	}

	t.Log(pp.Sprint(probData))
}

var theProb problem.Problem

func LoadProblem(t *testing.T) {
	var err error
	theProb, err = problem.LoadDir(probDataDir)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("statement: ", string(theProb.Stmt("zh")))
	t.Log("submission", pp.Sprint(theProb.SubmFields()))
}

var problemDumpDir string

func DumpProblem(t *testing.T) {
	err := theProb.DumpFile(path.Join(problemDumpDir, "dump.zip"))
	if err != nil {
		t.Error(err)
		return
	}
}

func ExtractProblem(t *testing.T) {
	dir := t.TempDir()
	_, err := problem.LoadDump(path.Join(problemDumpDir, "dump.zip"), dir)
	if err != nil {
		t.Error(err)
		return
	}

}
