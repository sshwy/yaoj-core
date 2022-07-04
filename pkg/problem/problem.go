package problem

import (
	"archive/zip"
	"encoding/json"
	"io"
	"log"
	"os"
	"path"
	"strconv"

	"github.com/sshwy/yaoj-core/pkg/utils"
	"golang.org/x/text/language"
)

type Problem interface {
	// 将题目打包为一个文件（压缩包）
	DumpFile(filename string) error
	// 获取题面，lang 见 http://www.lingoes.net/zh/translator/langcode.htm
	Stmt(lang string) []byte
	// 题解
	Tutr(lang string) []byte
	// 附加文件
	Assert(filename string) (*os.File, error)
	// 获取提交格式的数据表格
	SubmConf() SubmConf
	// 评测用的
	Data() *ProbData
	// 展示数据
	DataInfo() DataInfo
}

type DataInfo struct {
	IsSubtask  bool
	Fullscore  float64
	CalcMethod CalcMethod //计分方式
	Subtasks   []SubtaskInfo
	// 静态文件
	Static map[string]string //other properties of data
}

type SubtaskInfo struct {
	Id        int
	Fullscore float64
	Field     map[string]string //other properties of subtasks
	Tests     []TestInfo
}

type TestInfo struct {
	Id    int
	Field map[string]string //other properties of tests, i.e. in/output file path
}

type prob struct {
	data *ProbData
}

// 将题目打包为一个文件（压缩包）
func (r *prob) DumpFile(filename string) error {
	return zipDir(r.data.dir, filename)
}

func (r *prob) tryReadFile(filename string) []byte {
	ctnt, _ := os.ReadFile(path.Join(r.data.dir, filename))
	return ctnt
}

func (r *prob) Stmt(lang string) []byte {
	lang = GuessLang(lang)
	logger.Printf("Get statement lang=%s", lang)
	filename := r.data.Statement["s."+lang]
	return r.tryReadFile(filename)
}

func (r *prob) Tutr(lang string) []byte {
	lang = GuessLang(lang)
	logger.Printf("Get tutorial lang=%s", lang)
	filename := r.data.Statement["t."+lang]
	return r.tryReadFile(filename)
}

func (r *prob) Assert(filename string) (*os.File, error) {
	return os.Open(path.Join(r.data.dir, r.data.Statement[filename]))
}

// 获取提交格式的数据表格
func (r *prob) SubmConf() SubmConf {
	return r.data.Submission
}

func (r *prob) DataInfo() DataInfo {
	var res = DataInfo{
		IsSubtask:  r.data.IsSubtask(),
		Fullscore:  r.data.Fullscore,
		CalcMethod: r.data.CalcMethod,
		Static:     r.data.Static,
		Subtasks:   []SubtaskInfo{},
	}
	if res.IsSubtask {
		for i, task := range r.data.Subtasks.Record {
			var tests = []TestInfo{}
			for j, test := range r.data.Tests.Record {
				if test["_subtaskid"] != task["_subtaskid"] {
					continue
				}
				tests = append(tests, TestInfo{
					Id:    j,
					Field: copyRecord(test),
				})
			}

			score, _ := strconv.ParseFloat(task["_score"], 64)
			res.Subtasks = append(res.Subtasks, SubtaskInfo{
				Id:        i,
				Fullscore: score,
				Field:     task,
				Tests:     tests,
			})
		}
	} else {
		var tests = []TestInfo{}
		for j, test := range r.data.Tests.Record {
			tests = append(tests, TestInfo{
				Id:    j,
				Field: copyRecord(test),
			})
		}

		res.Subtasks = append(res.Subtasks, SubtaskInfo{
			Fullscore: r.data.Fullscore,
			Tests:     tests,
		})
	}
	return res
}

var _ Problem = (*prob)(nil)

// 加载一个题目文件夹
func LoadDir(dir string) (Problem, error) {
	data, err := LoadProbData(dir)
	if err != nil {
		return nil, err
	}
	return &prob{data: data}, nil
}

// 将打包的题目在空的文件夹下加载
func LoadDump(filename string, dir string) (Problem, error) {
	err := unzipSource(filename, dir)
	if err != nil {
		return nil, err
	}
	return LoadDir(dir)
}

var SupportLangs = []language.Tag{
	language.Chinese,
	language.English,
	language.Und,
}

var langMatcher = language.NewMatcher(SupportLangs)

var logger = log.New(os.Stderr, "[problem] ", log.LstdFlags|log.Lshortfile|log.Lmsgprefix)

func GuessLang(lang string) string {
	tag, _, _ := langMatcher.Match(language.Make(lang))
	if tag == language.Und {
		tag = SupportLangs[0]
	}
	base, _ := tag.Base()
	return base.String()
}

func (r *prob) Data() *ProbData {
	return r.data
}

// limitation for any file submitted
type SubmLimit struct {
	// 接受的语言，nil 表示所有语言
	Langs []utils.LangTag
	// 接受哪些类型的文件，必须设置值
	Accepted utils.CtntType
	// 文件大小，单位 byte
	Length uint32
}

// 存储文件的路径
type Submission map[string]string

// 加入提交文件
func (r Submission) Set(field string, name string) {
	r[field] = name
}

// 打包
func (r Submission) DumpFile(name string) error {
	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()

	w := zip.NewWriter(file)
	defer w.Close()

	var pathmap = map[string]string{}

	for field, name := range r {
		file, err := os.Open(name)
		if err != nil {
			return err
		}

		filename := field + "-" + path.Base(name)
		f, err := w.Create(filename)
		if err != nil {
			return err
		}

		_, err = io.Copy(f, file)
		if err != nil {
			return err
		}

		file.Close()

		pathmap[field] = filename
	}

	conf, err := w.Create("_config.json")
	if err != nil {
		return err
	}

	jsondata, err := json.Marshal(pathmap)
	if err != nil {
		return err
	}

	conf.Write(jsondata)
	return nil
}

// 解压
func LoadSubm(name string, dir string) (Submission, error) {
	err := unzipSource(name, dir)
	if err != nil {
		return nil, err
	}
	bconf, err := os.ReadFile(path.Join(dir, "_config.json"))
	if err != nil {
		return nil, err
	}
	var pathmap map[string]string
	if err := json.Unmarshal(bconf, &pathmap); err != nil {
		return nil, err
	}
	var res = Submission{}
	for field, name := range pathmap {
		res[field] = path.Join(dir, name)
	}
	return res, nil
}

// 提交文件配置
type SubmConf map[string]SubmLimit
