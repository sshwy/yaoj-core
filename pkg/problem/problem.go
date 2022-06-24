package problem

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"

	"github.com/sshwy/yaoj-core/pkg/utils"
	"github.com/sshwy/yaoj-core/pkg/workflow"
)

type Problem struct {
	Fullscore float64
	dir       string
	workflow  workflow.Workflow
	// "testcase" _subtaskid, _score ("average", {number})
	Tests table
	// "subtask" _subtaskid, _score
	Subtasks table
	// "static"
	Static table
	// "submission"
	Submission table
}

// Add file to r.dir/patch and return relative path
func (r *Problem) AddFile(name string, pathname string) (string, error) {
	if _, err := utils.CopyFile(pathname, path.Join(r.dir, "patch", name)); err != nil {
		return "", err
	}
	return path.Join("patch", name), nil
}

// export the problem's data to another empty dir and change itself to the new one
func (r *Problem) Export(dir string) error {
	os.Mkdir(path.Join(dir, "workflow"), os.ModePerm)
	graph_json, err := json.Marshal(r.workflow.WorkflowGraph)
	if err != nil {
		return err
	}
	if err := os.WriteFile(path.Join(dir, "workflow", "graph.json"), graph_json, 0644); err != nil {
		return err
	}
	os.Mkdir(path.Join(dir, "data"), os.ModePerm)
	os.Mkdir(path.Join(dir, "data", "tests"), os.ModePerm)
	os.Mkdir(path.Join(dir, "data", "subtasks"), os.ModePerm)
	os.Mkdir(path.Join(dir, "data", "static"), os.ModePerm)
	os.Mkdir(path.Join(dir, "patch"), os.ModePerm)
	var tests, subtasks, static table
	if tests, err = r.exportTable(r.Tests, dir, path.Join("data", "tests")); err != nil {
		return err
	}
	if subtasks, err = r.exportTable(r.Subtasks, dir, path.Join("data", "subtasks")); err != nil {
		return err
	}
	if static, err = r.exportTable(r.Static, dir, path.Join("data", "static")); err != nil {
		return err
	}
	r.Tests = tests
	r.Subtasks = subtasks
	r.Static = static

	prob_json, err := json.Marshal(*r)
	if err != nil {
		return err
	}
	if err := os.WriteFile(path.Join(dir, "problem.json"), prob_json, 0644); err != nil {
		return err
	}
	r.dir = dir
	return nil
}

func (r *Problem) exportTable(tb table, dir, dirtb string) (table, error) {
	log.Printf("exportTable %s", dirtb)
	var res table
	if res_json, err := json.Marshal(tb); err != nil {
		return tb, err
	} else {
		if err := json.Unmarshal(res_json, &res); err != nil {
			return tb, err
		}
	}

	// pp.Print(tb)
	for i, r2 := range res.Record {
		for k, v := range r2 {
			if k[0] == '_' { // internal field
				continue
			}
			name := fmt.Sprintf("%s%d%s", k, i, path.Ext(v))
			if _, err := utils.CopyFile(path.Join(r.dir, v), path.Join(dir, dirtb, name)); err != nil {
				return tb, err
			}
			r2[k] = path.Join(dirtb, name)
		}
	}
	// pp.Print(res, tb)
	return res, nil
}

func (r *Problem) SetWkflGraph(serial []byte) error {
	graph, err := workflow.Load(serial)
	if err != nil {
		return err
	}
	r.workflow.WorkflowGraph = graph
	return nil
}

// load problem from a dir
func Load(dir string) (*Problem, error) {
	serial, err := os.ReadFile(path.Join(dir, "problem.json"))
	if err != nil {
		return nil, err
	}
	var prob Problem
	if err := json.Unmarshal(serial, &prob); err != nil {
		return nil, err
	}
	// initialize
	prob.dir = dir
	wkgh, err := workflow.LoadFile(path.Join(dir, "workflow", "graph.json"))
	if err != nil {
		return nil, err
	}
	prob.workflow = workflow.Workflow{
		WorkflowGraph: wkgh,
		Analyzer:      workflow.DefaultAnalyzer{},
	}
	return &prob, nil
}

// create a new problem in an empty dir
func New(dir string) (*Problem, error) {
	var prob = Problem{
		dir: dir,
		workflow: workflow.Workflow{
			WorkflowGraph: &workflow.WorkflowGraph{},
			Analyzer:      workflow.DefaultAnalyzer{},
		},
		Tests:      newTable(),
		Subtasks:   newTable(),
		Static:     newTable(),
		Submission: newTable(),
	}
	if err := prob.Export(dir); err != nil {
		return nil, err
	}
	return &prob, nil
}

// Whether subtask is enabled.
func (r *Problem) IsSubtask() bool {
	return len(r.Subtasks.Field) > 0 && len(r.Subtasks.Record) > 0
}
func (r *Problem) toPathMap(rcd record) *map[string]string {
	res := map[string]string{}
	for k, v := range rcd {
		res[k] = path.Join(r.dir, v)
	}
	return &res
}

// Run all testcase in the dir.
func (r *Problem) Run(dir string, submission map[string]string) (*Result, error) {
	logger.Printf("run dir=%s", dir)
	// check submission
	for k := range r.Submission.Field {
		if _, ok := submission[k]; !ok {
			return nil, fmt.Errorf("submission missing field %s", k)
		}
	}

	var inboundPath = map[string]*map[string]string{
		"submission": (*map[string]string)(&submission),
	}
	if len(r.Static.Record) > 0 {
		inboundPath["static"] = r.toPathMap(r.Static.Record[0])
	}
	var result = Result{
		IsSubtask: r.IsSubtask(),
		Subtask:   []SubtResult{},
	}
	if r.IsSubtask() {
		for _, subtask := range r.Subtasks.Record {
			sub_res := SubtResult{
				Subtaskid: subtask["_subtaskid"],
				Testcase:  []workflow.Result{},
			}
			inboundPath["subtask"] = r.toPathMap(subtask)
			tests := r.testcaseOf(subtask["_subtaskid"])
			score, err := strconv.ParseFloat(subtask["_score"], 64)
			if err != nil {
				return nil, err
			}
			for _, test := range tests {
				inboundPath["testcase"] = r.toPathMap(test)
				res, err := workflow.Run(r.workflow, dir, inboundPath, score/float64(len(tests)))
				if err != nil {
					return nil, err
				}
				sub_res.Testcase = append(sub_res.Testcase, *res)
			}
			result.Subtask = append(result.Subtask, sub_res)
		}
	} else {
		sub_res := SubtResult{
			Testcase: []workflow.Result{},
		}
		for _, test := range r.Tests.Record {
			inboundPath["testcase"] = r.toPathMap(test)

			score := r.Fullscore / float64(len(r.Tests.Record))
			if f, err := strconv.ParseFloat(test["_score"], 64); err == nil {
				score = f
			}
			res, err := workflow.Run(r.workflow, dir, inboundPath, score)
			if err != nil {
				return nil, err
			}
			sub_res.Testcase = append(sub_res.Testcase, *res)
		}
		result.Subtask = append(result.Subtask, sub_res)
	}
	return &result, nil
}

func (r *Problem) testcaseOf(subtaskid string) []record {
	res := []record{}
	for _, test := range r.Tests.Record {
		if test["_subtaskid"] == subtaskid {
			res = append(res, test)
		}
	}
	return res
}

type Result struct {
	IsSubtask bool
	Subtask   []SubtResult
}
type SubtResult struct {
	Subtaskid string
	Testcase  []workflow.Result
}

var logger = log.New(os.Stderr, "[problem] ", log.LstdFlags|log.Lshortfile|log.Lmsgprefix)
