package problem

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/sshwy/yaoj-core/pkg/utils"
	"github.com/sshwy/yaoj-core/pkg/workflow"
)

// 子任务中测试点得分的汇总方式
type CalcMethod int

const (
	// default
	Mmin CalcMethod = iota
	Mmax
	Msum
)

// Problem result
type Result struct {
	IsSubtask  bool
	CalcMethod CalcMethod
	Subtask    []SubtResult
}

func (r Result) Byte() []byte {
	data, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	return data
}

// Subtask result
type SubtResult struct {
	Subtaskid string
	Fullscore float64
	Testcase  []workflow.Result
}

// Problem data module
type ProbData struct {
	// Usually 100.
	// Full score can be used to determine the point of testcase
	Fullscore  float64
	CalcMethod CalcMethod
	workflow   workflow.Workflow
	// "tests" _subtaskid, _score ("average", {number})
	Tests table
	// "subtask" _subtaskid, _score, _depend (separated by ",")
	Subtasks table
	// "submission" configuration
	Submission map[string]SubmLimit
	// "static"
	Static record
	// "statement"
	// Statement has 1 record. "s.{lang}", "t.{lang}" represents statement and
	// tutorial respectively. "_tl" "_ml" "_ol" denotes cpu time limit (ms),
	// real memory limit (MB) and output limit (MB) respectively.
	// Others are just filename.
	Statement record
	dir       string
}

// Add file to r.dir/patch and return relative path
func (r *ProbData) AddFile(name string, pathname string) (string, error) {
	name = path.Join("patch", name)
	logger.Printf("AddFile: %#v => %#v", pathname, name)
	if _, err := utils.CopyFile(pathname, path.Join(r.dir, name)); err != nil {
		return "", err
	}
	return name, nil
}

func (r *ProbData) AddFileReader(name string, file io.Reader) (string, error) {
	name = path.Join("patch", name)
	logger.Printf("AddFile: reader => %#v", name)
	destination, err := os.Create(path.Join(r.dir, name))
	if err != nil {
		return "", err
	}
	defer destination.Close()
	_, err = io.Copy(destination, file)
	if err != nil {
		return "", err
	}
	return name, nil
}

// export the problem's data to another empty dir and change itself to the new one
func (r *ProbData) Export(dir string) error {
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

	var tests, subtasks table
	var statement, static record
	if tests, err = r.exportTable(r.Tests, dir, path.Join("data", "tests")); err != nil {
		return err
	}
	if subtasks, err = r.exportTable(r.Subtasks, dir, path.Join("data", "subtasks")); err != nil {
		return err
	}
	if static, err = r.exportRecord(0, r.Static, dir, path.Join("data", "static")); err != nil {
		return err
	}
	if statement, err = r.exportRecord(0, r.Statement, dir, path.Join("statement")); err != nil {
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

func copyTable(tb table) (res table) {
	if res_json, err := json.Marshal(tb); err != nil {
		panic(err)
	} else {
		if err := json.Unmarshal(res_json, &res); err != nil {
			panic(err)
		}
	}
	return
}

func copyRecord(rcd record) (res record) {
	res = record{}
	for field, val := range rcd {
		res[field] = val
	}
	return
}

func (r *ProbData) exportRecord(id int, rcd record, newroot, dircd string) (res record, err error) {
	log.Printf("Export Record #%d %#v", id, dircd)
	res = make(record)
	for field, val := range rcd {
		if field[0] == '_' { // private field
			res[field] = rcd[field]
		} else {
			name := fmt.Sprintf("%s%d%s", field, id, path.Ext(val))
			if _, err := utils.CopyFile(path.Join(r.dir, val), path.Join(newroot, dircd, name)); err != nil {
				return res, err
			}
			res[field] = path.Join(dircd, name)
		}
	}
	return res, nil
}

func (r *ProbData) exportTable(tb table, newroot, dirtb string) (table, error) {
	log.Printf("Export Table %#v", dirtb)
	res := copyTable(tb)

	for i, record := range tb.Record {
		for field := range record {
			if _, ok := tb.Field[field]; !ok {
				return tb, fmt.Errorf("invalid field %s in record #%d", field, i)
			}
		}
		rcd, err := r.exportRecord(i, record, newroot, dirtb)
		if err != nil {
			return tb, err
		}
		res.Record[i] = rcd
	}
	return res, nil
}

// Set workflow graph
func (r *ProbData) SetWkflGraph(serial []byte) error {
	graph, err := workflow.Load(serial)
	if err != nil {
		return err
	}
	r.workflow.WorkflowGraph = graph
	return nil
}

// load problem from a dir
func LoadProbData(dir string) (*ProbData, error) {
	serial, err := os.ReadFile(path.Join(dir, "problem.json"))
	if err != nil {
		return nil, err
	}
	var prob ProbData
	if err := json.Unmarshal(serial, &prob); err != nil {
		return nil, err
	}
	// initialize
	absdir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	prob.dir = absdir
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
func NewProbData(dir string) (*ProbData, error) {
	graph := workflow.NewGraph()
	var prob = ProbData{
		dir: dir,
		workflow: workflow.Workflow{
			WorkflowGraph: &graph,
			Analyzer:      workflow.DefaultAnalyzer{},
		},
		Tests:      newTable(),
		Subtasks:   newTable(),
		Static:     make(record),
		Submission: map[string]SubmLimit{},
		Statement:  make(record),
	}
	if err := prob.Export(dir); err != nil {
		return nil, err
	}
	return &prob, nil
}

// Whether subtask is enabled.
func (r *ProbData) IsSubtask() bool {
	return len(r.Subtasks.Field) > 0 && len(r.Subtasks.Record) > 0
}

// get the workflow
func (r *ProbData) Workflow() workflow.Workflow {
	return r.workflow
}

// get problem dir
func (r *ProbData) Dir() string {
	return r.dir
}

// Set statement content to file in r.dir
func (r *ProbData) SetStmt(lang string, file string) {
	r.Statement["s."+GuessLang(lang)] = file
}

func (r *ProbData) SetValFile(rcd record, field string, filename string) error {
	pin, err := r.AddFile(utils.RandomString(5)+"_"+path.Base(filename), filename)
	if err != nil {
		return err
	}
	rcd[field] = pin
	return nil
}

// remove problem (dir)
func (r *ProbData) Finalize() error {
	logger.Printf("finalize %q", r.dir)
	return os.RemoveAll(r.dir)
}
