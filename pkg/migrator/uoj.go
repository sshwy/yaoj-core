package migrator

import (
	"bytes"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/sshwy/yaoj-core/pkg/problem"
	"github.com/sshwy/yaoj-core/pkg/utils"
	"github.com/sshwy/yaoj-core/pkg/workflow"
)

// 格式：
//
//   prob/
//     data/ # 里面放数据（就是svn的1文件夹里的内容。也就是说data/problem.conf是配置文件
//     statement/ # 放public资源，pdf啥的。不建子文件夹
//       statement/statement.md # 特殊，放题面内容（UOJ好像是没有多语言题面的，所以默认中文
//       statement/tutorial.pdf # 题解
//
type Uoj struct{}

var _ Migrator = Uoj{}

func (r Uoj) Migrate(src string, dest string) (Problem, error) {
	tmpdir := path.Join(os.TempDir(), "uoj-yaoj-migrator")
	os.RemoveAll(tmpdir)
	if err := os.MkdirAll(tmpdir, os.ModePerm); err != nil {
		return nil, err
	}
	prob, err := problem.NewProbData(tmpdir)
	if err != nil {
		return nil, err
	}
	prob.Fullscore = 100

	// parse statement (ignore error)
	filepath.Walk(path.Join(src, "statement"), func(pathname string, info fs.FileInfo, err error) error {
		if err != nil {
			logger.Printf("prevent panic by handling failure accessing a path %q: %v", pathname, err)
			return err
		}
		if info.IsDir() && info.Name() != "statement" {
			logger.Printf("skipping a dir: %#v", info.Name())
			return filepath.SkipDir
		}

		if info.IsDir() {
			return nil
		}

		// logger.Printf("Access file: %q", pathname)
		basename := path.Base(pathname)
		patch, err := prob.AddFile(basename, pathname)
		if err != nil {
			return err
		}
		prob.Statement[basename] = patch
		return nil
	})
	if _, ok := prob.Statement["statement.md"]; ok {
		prob.SetStmt("zh", prob.Statement["statement.md"])
	}

	fconf, err := os.ReadFile(path.Join(src, "data", "problem.conf"))
	if err != nil {
		return nil, err
	}
	conf := parseConf(fconf)

	if conf["use_builtin_judger"] != "on" {
		panic("gg")
	}

	// parse tests
	prob.Tests.Fields().Add("input")
	prob.Tests.Fields().Add("output")
	prob.Tests.Fields().Add("_score")

	n_tests := parseInt(conf["n_tests"])
	for i := 1; i <= n_tests; i++ {
		input := fmt.Sprint(conf["input_pre"], i, ".", conf["input_suf"])
		output := fmt.Sprint(conf["output_pre"], i, ".", conf["output_suf"])

		rcd := prob.Tests.Records().New()
		err = prob.SetValFile(rcd, "input", path.Join(src, "data", input))
		if err != nil {
			return nil, err
		}
		err = prob.SetValFile(rcd, "output", path.Join(src, "data", output))
		if err != nil {
			return nil, err
		}
	}

	// parse checker
	if _, ok := conf["use_builtin_checker"]; ok {
		logger.Printf("use builtin checker: %q", conf["use_builtin_checker"])
		// copy checker
		file, err := asserts.Open(path.Join("asserts", "checker", conf["use_builtin_checker"]+".cpp"))
		if err != nil {
			return nil, err
		}
		pchk, err := prob.AddFileReader("checker_uoj.cpp", file)
		if err != nil {
			return nil, err
		}
		file.Close()
		prob.Static["checker"] = pchk
	} else { // custom checker
		pchk, err := prob.AddFile("checker_custom.cpp", path.Join(src, "data", "chk.cpp"))
		if err != nil {
			return nil, err
		}
		prob.Static["checker"] = pchk
	}

	// parse limitation
	tl := parseInt(conf["time_limit"])
	ml := parseInt(conf["memory_limit"])
	ol := parseInt(conf["output_limit"])
	limReader := bytes.NewReader([]byte(fmt.Sprintf(
		"%d %d %d %d %d %d %d",
		1000*60, // 1min
		1000*tl,
		0, // no limit
		1024*1024*ml,
		1024*1024*ml,
		1024*1024*ol,
		50,
	)))
	plim, err := prob.AddFileReader("lim.txt", limReader)
	if err != nil {
		return nil, err
	}
	prob.Static["limitation"] = plim
	prob.Statement["_tl"] = fmt.Sprint(tl * 1000)
	prob.Statement["_ml"] = conf["memory_limit"]
	prob.Statement["_ol"] = conf["output_limit"]

	var builder workflow.Builder
	builder.SetNode("compile_source", "compiler:auto", false)
	builder.SetNode("compile_checker", "compiler:testlib", false)
	builder.SetNode("check", "checker:testlib", false)
	builder.SetNode("run", "runner:stdio", true)
	builder.AddInbound(workflow.Gstatic, "limitation", "run", "limit")
	builder.AddInbound(workflow.Gstatic, "checker", "compile_checker", "source")
	builder.AddInbound(workflow.Gsubm, "source", "compile_source", "source")
	builder.AddInbound(workflow.Gtests, "input", "run", "stdin")
	builder.AddInbound(workflow.Gtests, "input", "check", "input")
	builder.AddInbound(workflow.Gtests, "output", "check", "answer")
	builder.AddEdge("compile_source", "result", "run", "executable")
	builder.AddEdge("compile_checker", "result", "check", "checker")
	builder.AddEdge("run", "stdout", "check", "output")
	graph, err := builder.WorkflowGraph()
	if err != nil {
		return nil, err
	}
	err = prob.SetWkflGraph(graph.Serialize())
	if err != nil {
		return nil, err
	}

	// parse subtask
	if sNsubt, ok := conf["n_subtasks"]; ok {
		nsubt, _ := strconv.ParseInt(sNsubt, 10, 32)
		prob.Tests.Fields().Add("_subtaskid")
		prob.Subtasks.Fields().Add("_subtaskid")
		prob.Subtasks.Fields().Add("_score")

		logger.Printf("nsubtask = %d", nsubt)

		las := 0
		for i := 1; i <= int(nsubt); i++ {
			endid, _ := strconv.ParseInt(conf[fmt.Sprint("subtask_end_", i)], 10, 32)
			score, _ := strconv.ParseInt(conf[fmt.Sprint("subtask_score_", i)], 10, 32)
			record := prob.Subtasks.Records().New()
			record["_subtaskid"] = fmt.Sprint("subtask_", i)
			record["_score"] = fmt.Sprint(score)

			for j := las; j < int(endid); j++ {
				prob.Tests.Record[j]["_subtaskid"] = record["_subtaskid"]
			}
			las = int(endid)
		}
	} else {
		prob.CalcMethod = problem.Msum
	}

	// analyzer
	// panic("not complete")
	prob.Submission["source"] = problem.SubmLimit{
		Length:   1024 * 64,
		Accepted: utils.Csource,
	}

	os.MkdirAll(dest, os.ModePerm)
	err = prob.Export(dest)
	if err != nil {
		return nil, err
	}
	return problem.LoadDir(dest)
}

// 对于每一行，解析前两个不含空格的字符串分别作为字段和值
func parseConf(content []byte) (res map[string]string) {
	res = map[string]string{}
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.Trim(line, " \t\f\n\r")
		if line == "" {
			continue
		}
		tokens := strings.Split(line, " ")
		var directive, val string
		for _, token := range tokens {
			if token == "" {
				continue
			}
			if directive == "" {
				directive = token
			} else {
				val = token
				break
			}
		}
		if directive == "" {
			panic(fmt.Sprintf("invalid line %#v", line))
		}
		res[directive] = val
	}
	logger.Printf("conf: %+v", res)
	return
}

func parseInt(s string) int {
	res, _ := strconv.ParseInt(s, 10, 32)
	return int(res)
}

var logger = log.New(os.Stderr, "[migrator] ", log.LstdFlags|log.Lshortfile|log.Lmsgprefix)
