package problem

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

// Problem datagroup consists of a directory (datagroup's name) containing a
// series of files without subdirectories. Each file is human-readbly named
//
//     [record rank].[field].[arbtrary suffix]
//
type ProbDtgp struct {
	dir    string
	fields map[string]bool
	// suffix
	records []*map[string]string
}

func (r *ProbDtgp) rcpath(id int, field string, suf string) string {
	return path.Join(r.dir, fmt.Sprint(id, ".", field, ".", suf))
}
func (r *ProbDtgp) Name() string {
	return path.Base(r.dir)
}

// a list of unique strings denoting all fields, sorted by alphabet
func (r *ProbDtgp) Fields() []string {
	res := []string{}
	for k := range r.fields {
		res = append(res, k)
	}
	return res
}

// a table of all records.
func (r *ProbDtgp) Records() []*map[string]string {
	return r.records
}

func (r *ProbDtgp) AddField(name string) error {
	logger.Printf(`AddField "%s" for group %s`, name, r.dir)
	if _, ok := r.fields[name]; ok {
		return fmt.Errorf("field \"%s\" has already existed", name)
	}
	r.fields[name] = true
	suf := "txt"
	for i, record := range r.records {
		(*record)[name] = suf // default suffix
		file, err := os.Create(r.rcpath(i, name, suf))
		if err != nil {
			return err
		}
		file.Close()
	}
	return nil
}

func (r *ProbDtgp) RemoveField(name string) error {
	logger.Printf(`RemoveField "%s" for group %s`, name, r.dir)
	if _, ok := r.fields[name]; !ok {
		return fmt.Errorf("field \"%s\" not exist", name)
	}
	delete(r.fields, name)
	for i, record := range r.records {
		os.Remove(r.rcpath(i, name, (*record)[name]))
		delete(*record, name)
	}
	return nil
}

// append an empty record
func (r *ProbDtgp) NewRecord() error {
	logger.Printf(`NewRecord for group %s`, r.dir)
	r.records = append(r.records, &map[string]string{})
	id := len(r.records) - 1
	record := *r.records[id]
	for k := range r.fields {
		suf := "txt"
		record[k] = suf // default suffix
		file, err := os.Create(r.rcpath(id, k, suf))
		if err != nil {
			return err
		}
		file.Close()
	}
	return nil
}

// id start from 0.
func (r *ProbDtgp) RemoveRecord(id int) error {
	logger.Printf(`RemoveRecord %d for group %s`, id, r.dir)
	if id < 0 || id > len(r.records) {
		return fmt.Errorf("invalid id")
	}
	record := *r.records[id]
	for k, suf := range record {
		os.Remove(r.rcpath(id, k, suf))
	}
	r.records = append(r.records[:id], r.records[id+1:]...)
	return nil
}

func (r *ProbDtgp) AlterValue(id int, field string, filepath string) error {
	logger.Printf(`AlterValue id=%d field=%s filepath=%s for group %s`, id, field, filepath, r.dir)
	ext := path.Ext(filepath)
	if ext == "" {
		return fmt.Errorf("invalid filepath: no suffix found")
	}
	ext = ext[1:] // remove dot
	fin, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer fin.Close()

	record := *r.records[id]
	os.Remove(r.rcpath(id, field, record[field]))
	record[field] = ext

	file, err := os.Create(r.rcpath(id, field, ext))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, fin)
	if err != nil {
		return err
	}
	return nil
}

func LoadDtgp(dirdtgp string) (*ProbDtgp, error) {
	var group ProbDtgp = ProbDtgp{
		dir:     dirdtgp,
		fields:  map[string]bool{},
		records: []*map[string]string{},
	}
	var fields = map[string]bool{}
	var recordmap = map[int64]*map[string]string{}
	err := filepath.WalkDir(dirdtgp, func(pathname string, d fs.DirEntry, _ error) error {
		if d.IsDir() {
			if pathname == dirdtgp { // root
				return nil
			}
			return fs.SkipDir
		}
		name := path.Base(d.Name())
		parts := strings.Split(name, ".")
		logger.Print(parts)
		if len(parts) != 3 {
			return fmt.Errorf("invalid filename %s", name)
		}
		id, err := strconv.ParseInt(parts[0], 10, 32)
		if err != nil {
			return err
		}
		if recordmap[id] == nil {
			recordmap[id] = &map[string]string{}
		}
		if v, ok := (*recordmap[id])[parts[1]]; ok {
			return fmt.Errorf("record[%d][%s] conflict (%s, %s)", id, parts[1], parts[2], v)
		}
		(*recordmap[id])[parts[1]] = parts[2]
		fields[parts[1]] = true
		return nil
	})
	if err != nil {
		return nil, err
	}
	recordcnt := len(recordmap)
	for i := 0; i < recordcnt; i++ {
		if _, ok := recordmap[int64(i)]; !ok {
			return nil, fmt.Errorf("invalid record: missing id %d", i)
		}
		record := recordmap[int64(i)]
		if len(*record) != len(fields) {
			return nil, fmt.Errorf("invalid record: incorrect field set (id=%d)", i)
		}
		for k := range fields {
			if _, ok := (*record)[k]; !ok {
				return nil, fmt.Errorf("invalid record: missing field %s (id=%d)", k, i)
			}
		}
		group.records = append(group.records, record)
	}
	group.fields = fields

	return &group, nil
}
