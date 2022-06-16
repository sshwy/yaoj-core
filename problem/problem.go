package problem

import (
	"fmt"
	"hash/crc64"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"

	"github.com/sshwy/yaoj-core/workflow"
)

var logger = log.New(os.Stderr, "[problem] ", log.LstdFlags|log.Lshortfile|log.Lmsgprefix)

// Root directory of a problem with some helper functions
type ProbRoot string

// string(r)
func (r ProbRoot) Root() string {
	return string(r)
}

// r/datagroup
func (r ProbRoot) DtgpDir() string {
	return path.Join(r.Root(), "datagroup")
}

// r/statement
func (r ProbRoot) StmtDir() string {
	return path.Join(r.Root(), "statement")
}

// r/statement/statement.md
func (r ProbRoot) StmtMkdn() string {
	return path.Join(r.StmtDir(), "statement.md")
}

// r/workflow
func (r ProbRoot) WkflDir() string {
	return path.Join(r.Root(), "workflow")
}

// r/workflow/graph.json
func (r ProbRoot) WkflGraph() string {
	return path.Join(r.WkflDir(), "graph.json")
}

// r/workflow/analyzer.go
func (r ProbRoot) WkflAnyz() string {
	return path.Join(r.WkflDir(), "analyzer.go")
}

// A testcase of a problem
type Testcase struct {
	// map[datagroup_name]record_id
	groups  map[string]int
	problem *Problem
}

func (r *Testcase) ID() uint64 {
	var a = make([]string, 0, len(r.groups))
	for name, id := range r.groups {
		a = append(a, fmt.Sprint(name, id))
	}
	sort.Slice(a, func(i, j int) bool {
		return a[i] < a[j]
	})
	hash := crc64.New(crc64.MakeTable(crc64.ISO))
	for _, v := range a {
		hash.Write([]byte(v))
	}
	return hash.Sum64()
}

func (r *Testcase) InboundPath(submission map[string]string) map[string]*map[string]string {
	a := map[string]*map[string]string{}
	for name, id := range r.groups {
		a[name] = r.problem.groups[name].Record(id).PathMap()
	}
	a["submission"] = &submission
	return a
}

// Run the testcase in dir with submission and fullscore provided.
func (r *Testcase) Run(dir string, submission map[string]string, fullscore float64) (*workflow.Result, error) {
	return workflow.Run(r.problem.workflow, dir, r.InboundPath(submission), fullscore)
}

type Problem struct {
	// where it store
	dir       ProbRoot
	statement []byte
	groups    map[string]*ProbDtgp
	workflow  workflow.Workflow
}

// Get problem statement.
func (r *Problem) Stmt() []byte {
	return r.statement
}

// Set problem statement.
func (r *Problem) SetStmt(content []byte) error {
	err := os.WriteFile(r.dir.StmtMkdn(), content, 0644)
	if err != nil {
		return err
	}
	r.statement = content
	return nil
}

// Create a new datagroup in dir/datagroup/[name]/
func (r *Problem) NewDataGroup(name string) (*ProbDtgp, error) {
	if _, ok := r.groups[name]; ok {
		return nil, fmt.Errorf("datagroup has already existed")
	}
	err := os.Mkdir(path.Join(r.dir.DtgpDir(), name), os.ModePerm)
	if err != nil {
		return nil, err
	}
	dtgp, err := LoadDtgp(path.Join(r.dir.DtgpDir(), name))
	if err != nil {
		return nil, err
	}
	r.groups[name] = dtgp
	return dtgp, nil
}

// Get a datagroup in dir/datagroup/[name]/, nil if not found
func (r *Problem) DataGroup(name string) *ProbDtgp {
	return r.groups[name]
}

func (r *Problem) SetWkflGraph(serial []byte) error {
	graph, err := workflow.Load(serial)
	if err != nil {
		return err
	}
	r.workflow.WorkflowGraph = graph
	return nil
}

// Enumerate all records in all datagroups to generate all testcases.
func (r *Problem) Testcase() []Testcase {
	dup := func(m map[string]int) map[string]int {
		m2 := make(map[string]int)
		for k, v := range m {
			m2[k] = v
		}
		return m2
	}

	groups := []map[string]int{{}}

	for name, dtgp := range r.groups {
		logger.Print(name)
		pathmap := dtgp.Record(0).PathMap()
		if pathmap == nil {
			panic("empty datagroup")
		}
		for i := range groups {
			groups[i][name] = 0
		}

		cnt := len(groups)

		for i := 1; i < dtgp.Len(); i++ {
			var more = make([]map[string]int, cnt)
			for j := 0; j < cnt; j++ {
				more[j] = dup(groups[j])
				more[j][name] = i
			}
			groups = append(groups, more...)
		}
	}

	var testcases []Testcase = make([]Testcase, 0, len(groups))
	// pp.Print(groups)
	for _, v := range groups {
		testcases = append(testcases, Testcase{
			groups:  v,
			problem: r,
		})
	}
	return testcases
}

// Load a problem from dir.
func Load(dir string) (*Problem, error) {
	root := ProbRoot(dir)
	statement, _ := os.ReadFile(root.StmtMkdn())

	group := map[string]*ProbDtgp{}
	err := filepath.WalkDir(root.DtgpDir(), func(pathname string, d fs.DirEntry, _ error) error {
		if pathname == root.DtgpDir() {
			return nil
		}
		if !d.IsDir() {
			return nil
		}
		dtgp, err := LoadDtgp(pathname)
		if err != nil {
			return err
		}
		group[dtgp.Name()] = dtgp
		return fs.SkipDir
	})
	if err != nil {
		return nil, err
	}

	wkgh, err := workflow.LoadFile(root.WkflGraph())
	if err != nil {
		return nil, err
	}

	logger.Print("custom analyzer not loaded!")

	prob := Problem{
		dir:       root,
		statement: statement,
		groups:    group,
		workflow: workflow.Workflow{
			WorkflowGraph: wkgh,
			Analyzer:      workflow.DefaultAnalyzer{},
		},
	}
	return &prob, nil
}

// create a new problem in an empty dir
func New(dir string) (*Problem, error) {
	root := ProbRoot(dir)
	if err := os.Mkdir(root.StmtDir(), os.ModePerm); err != nil {
		return nil, err
	}
	if err := os.Mkdir(root.DtgpDir(), os.ModePerm); err != nil {
		return nil, err
	}
	if err := os.Mkdir(root.WkflDir(), os.ModePerm); err != nil {
		return nil, err
	}
	file, err := os.Create(root.StmtMkdn())
	if err != nil {
		return nil, err
	}
	file.Close()
	file, err = os.Create(root.WkflGraph())
	if err != nil {
		return nil, err
	}

	if _, err := file.Write((&workflow.WorkflowGraph{
		Node:    []workflow.Node{},
		Inbound: map[string]*map[string]workflow.Bound{},
	}).Serialize()); err != nil {
		return nil, err
	}
	file.Close()
	return Load(dir)
}
