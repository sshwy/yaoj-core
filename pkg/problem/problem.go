package problem

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

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
	SubmFields() []string
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
func (r *prob) SubmFields() (res []string) {
	res = []string{}
	for field := range r.data.Submission.Field {
		res = append(res, field)
	}
	return
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

func zipDir(root string, dest string) error {
	root = path.Clean(root)
	file, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer file.Close()

	w := zip.NewWriter(file)
	defer w.Close()

	walker := func(pathname string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		file, err := os.Open(pathname)
		if err != nil {
			return err
		}
		defer file.Close()

		if pathname[:len(root)] != root {
			return fmt.Errorf("invalid path %s", pathname)
		}
		zippath := pathname[len(root)+1:]

		// Ensure that `pathname` is not absolute; it should not start with "/".
		// This snippet happens to work because I don't use
		// absolute pathnames, but ensure your real-world code
		// transforms pathname into a zip-root relative pathname.
		logger.Printf("Create %#v\n", zippath)
		f, err := w.Create(zippath)
		if err != nil {
			return err
		}

		_, err = io.Copy(f, file)
		if err != nil {
			return err
		}

		return nil
	}
	err = filepath.Walk(root, walker)
	if err != nil {
		return err
	}
	return nil
}

// https://gosamples.dev/unzip-file/
func unzipSource(source, destination string) error {
	// 1. Open the zip file
	reader, err := zip.OpenReader(source)
	if err != nil {
		return err
	}
	defer reader.Close()

	// 2. Get the absolute destination path
	destination, err = filepath.Abs(destination)
	if err != nil {
		return err
	}

	// 3. Iterate over zip files inside the archive and unzip each of them
	for _, f := range reader.File {
		err := unzipFile(f, destination)
		if err != nil {
			return err
		}
	}

	return nil
}

func unzipFile(f *zip.File, destination string) error {
	// 4. Check if file paths are not vulnerable to Zip Slip
	filePath := filepath.Join(destination, f.Name)
	if !strings.HasPrefix(filePath, filepath.Clean(destination)+string(os.PathSeparator)) {
		return fmt.Errorf("invalid file path: %s", filePath)
	}

	// 5. Create directory tree
	if f.FileInfo().IsDir() {
		if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
			return err
		}
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return err
	}

	// 6. Create a destination file for unzipped content
	destinationFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	// 7. Unzip the content of a file and copy it to the destination file
	zippedFile, err := f.Open()
	if err != nil {
		return err
	}
	defer zippedFile.Close()

	if _, err := io.Copy(destinationFile, zippedFile); err != nil {
		return err
	}
	return nil
}

func GuessLang(lang string) string {
	tag, _, _ := langMatcher.Match(language.Make(lang))
	if tag == language.Und {
		tag = SupportLangs[0]
	}
	base, _ := tag.Base()
	return base.String()
}
