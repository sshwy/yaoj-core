package run

import (
	"fmt"
	"path"
	"strconv"

	"github.com/sshwy/yaoj-core/pkg/problem"
	"github.com/sshwy/yaoj-core/pkg/workflow"
)

// change a record's relative path to real path
func toPathMap(r *problem.ProbData, rcd map[string]string) *map[string]string {
	res := map[string]string{}
	for k, v := range rcd {
		res[k] = path.Join(r.Dir(), v)
	}
	return &res
}
func testcaseOf(r *problem.ProbData, subtaskid string) []map[string]string {
	res := []map[string]string{}
	for _, test := range r.Tests.Record {
		if test["_subtaskid"] == subtaskid {
			res = append(res, test)
		}
	}
	return res
}

// Run all testcase in the dir.
func RunProblem(r *problem.ProbData, dir string, submission map[string]string) (*problem.Result, error) {
	logger.Printf("run dir=%s", dir)
	// check submission
	for k := range r.Submission {
		if _, ok := submission[k]; !ok {
			return nil, fmt.Errorf("submission missing field %s", k)
		}
	}

	var inboundPath = map[workflow.Groupname]*map[string]string{
		workflow.Gsubm: (*map[string]string)(&submission),
	}
	inboundPath[workflow.Gstatic] = toPathMap(r, r.Static)
	var result = problem.Result{
		IsSubtask: r.IsSubtask(),
		Subtask:   []problem.SubtResult{},
	}
	if r.IsSubtask() {
		for _, subtask := range r.Subtasks.Record {
			sub_res := problem.SubtResult{
				Subtaskid: subtask["_subtaskid"],
				Testcase:  []workflow.Result{},
			}
			inboundPath[workflow.Gsubt] = toPathMap(r, subtask)
			tests := testcaseOf(r, subtask["_subtaskid"])
			score, err := strconv.ParseFloat(subtask["_score"], 64)
			if err != nil {
				return nil, err
			}
			sub_res.Fullscore = score
			for _, test := range tests {
				inboundPath[workflow.Gtests] = toPathMap(r, test)
				res, err := RunWorkflow(r.Workflow(), dir, inboundPath, score/float64(len(tests)))
				if err != nil {
					return nil, err
				}
				sub_res.Testcase = append(sub_res.Testcase, *res)
			}
			result.Subtask = append(result.Subtask, sub_res)
		}
	} else {
		sub_res := problem.SubtResult{
			Testcase: []workflow.Result{},
		}
		for _, test := range r.Tests.Record {
			inboundPath[workflow.Gtests] = toPathMap(r, test)

			score := r.Fullscore / float64(len(r.Tests.Record))
			if f, err := strconv.ParseFloat(test["_score"], 64); err == nil {
				score = f
			}
			res, err := RunWorkflow(r.Workflow(), dir, inboundPath, score)
			if err != nil {
				return nil, err
			}
			sub_res.Testcase = append(sub_res.Testcase, *res)
		}
		sub_res.Fullscore = r.Fullscore
		result.Subtask = append(result.Subtask, sub_res)
	}
	return &result, nil
}
