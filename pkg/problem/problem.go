package problem

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/sshwy/yaoj-core/pkg/utils"
	"github.com/sshwy/yaoj-core/pkg/workflow"
)

// Problem result
type Result struct {
	IsSubtask bool
	Subtask   []SubtResult
}

// Subtask result
type SubtResult struct {
	Subtaskid string
	Testcase  []workflow.Result
}

type Problem struct {
	// usually 100
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
	// "statement"
	Statement table
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
	os.Mkdir(path.Join(dir, "statement"), os.ModePerm)

	var tests, subtasks, static, statement table
	if tests, err = r.exportTable(r.Tests, dir, path.Join("data", "tests")); err != nil {
		return err
	}
	if subtasks, err = r.exportTable(r.Subtasks, dir, path.Join("data", "subtasks")); err != nil {
		return err
	}
	if static, err = r.exportTable(r.Static, dir, path.Join("data", "static")); err != nil {
		return err
	}
	if statement, err = r.exportTable(r.Statement, dir, path.Join("statement")); err != nil {
		return err
	}

	// modify r from now
	r.Tests = tests
	r.Subtasks = subtasks
	r.Static = static
	r.Statement = statement

	prob_json, err := json.Marshal(*r)
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile(path.Join(dir, "problem.json"), prob_json, 0644); err != nil {
		panic(err)
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
	for i, record := range res.Record {
		for field, val := range record {
			if field[0] == '_' { // private field
				continue
			}
			name := fmt.Sprintf("%s%d%s", field, i, path.Ext(val))
			if _, err := utils.CopyFile(path.Join(r.dir, val), path.Join(dir, dirtb, name)); err != nil {
				return tb, err
			}
			record[field] = path.Join(dirtb, name)
		}
	}
	// pp.Print(res, tb)
	return res, nil
}

// Set workflow graph
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
		Statement:  newTable(),
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

// get all testcase of a subtask
func (r *Problem) TestcaseOf(subtaskid string) []record {
	res := []record{}
	for _, test := range r.Tests.Record {
		if test["_subtaskid"] == subtaskid {
			res = append(res, test)
		}
	}
	return res
}

// get the workflow
func (r *Problem) Workflow() workflow.Workflow {
	return r.workflow
}

// get problem dir
func (r *Problem) Dir() string {
	return r.dir
}
