package problem_test

import (
	"path"
	"testing"

	"github.com/bitfield/script"
	"github.com/k0kubun/pp/v3"
	"github.com/sshwy/yaoj-core/problem"
)

func TestNew(t *testing.T) {
	dir := t.TempDir()
	prob, err := problem.New(dir)
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

	prob.Fullscore = 100
	prob.Tests.Fields().Add("input")
	prob.Tests.Fields().Add("answer")
	prob.Tests.Fields().Add("_subtaskid")
	prob.Tests.Fields().Add("_score")
	r0 := prob.Tests.Records().New()
	r0["input"], err = prob.AddFile("a.in", path.Join(dir, "a.in"))
	if err != nil {
		t.Error(err)
		return
	}
	r0["answer"], err = prob.AddFile("a.ans", path.Join(dir, "a.ans"))
	if err != nil {
		t.Error(err)
		return
	}
	r0["_score"] = "average"

	r1 := prob.Tests.Records().New()
	r1["input"] = r0["input"]
	r1["answer"] = r0["answer"]
	r1["_score"] = "average"

	r2 := prob.Tests.Records().New()
	r2["input"] = r0["input"]
	r2["answer"] = r0["answer"]
	r2["_score"] = "average"

	r3 := prob.Tests.Records().New()
	r3["input"] = r0["input"]
	r3["answer"] = r0["answer"]
	r3["_score"] = "average"

	prob.Static.Fields().Add("limitation")
	prob.Static.Fields().Add("compilescript")
	o0 := prob.Static.Records().New()
	o0["limitation"] = "cpl.txt"
	o0["compilescript"] = "script.sh"

	// net adjuestment
	err = prob.SetWkflGraph([]byte(`
{
    "Node": {
        "check": {
            "ProcName": "checker:hcmp"
        },
        "compile": {
            "ProcName": "compiler"
        },
        "run": {
            "ProcName": "runner:stdio",
            "Key": true
        }
    },
    "Edge": [
        {
            "From": {
                "Name": "compile",
                "LabelIndex": 0
            },
            "To": {
                "Name": "run",
                "LabelIndex": 0
            }
        },
        {
            "From": {
                "Name": "run",
                "LabelIndex": 0
            },
            "To": {
                "Name": "check",
                "LabelIndex": 0
            }
        }
    ],
    "Inbound": {
        "static": {
            "compilescript": [
                {
                    "Name": "compile",
                    "LabelIndex": 1
                }
            ],
            "limitation": [
                {
                    "Name": "run",
                    "LabelIndex": 2
                }
            ]
        },
        "submission": {
            "source": [
                {
                    "Name": "compile",
                    "LabelIndex": 0
                }
            ]
        },
        "testcase": {
            "answer": [
                {
                    "Name": "check",
                    "LabelIndex": 1
                }
            ],
            "input": [
                {
                    "Name": "run",
                    "LabelIndex": 1
                }
            ]
        }
    }
}
	`))
	if err != nil {
		t.Error(err)
		return
	}
	prob.Submission.Fields().Add("source")
	// pp.Print(prob)

	if err := prob.Export(t.TempDir()); err != nil {
		t.Error(err)
		return
	}

	res, err := prob.Run(t.TempDir(), map[string]string{
		"source": path.Join(dir, "src.cpp"),
	})
	if err != nil {
		t.Error(err)
		return
	}
	pp.Print(res)
}

// func TestLoad(t *testing.T) {
// 	_, err := problem.Load("testdata/prob")
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	// t.Log(pp.Sprint(prob))
// }
