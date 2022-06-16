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
// series of files without subdirectories. Each file is human-readbly named.
type ProbDtgp struct {
	dir    string
	fields map[string]bool
	// suffix
	records []*Record
}

func (r *ProbDtgp) rcpath(id int, field string, suf string) string {
	return path.Join(r.dir, fmt.Sprint(id, ".", field, ".", suf))
}

// Basically the last directory of the path
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

// Number of records.
func (r *ProbDtgp) Len() int {
	return len(r.records)
}

// a table of all records.
func (r *ProbDtgp) Records() []*Record {
	return r.records
}

// id start from 0.
func (r *ProbDtgp) Record(id int) *Record {
	if id < 0 || id >= len(r.records) {
		return nil
	}
	return r.records[id]
}

func (r *ProbDtgp) AddField(name string) error {
	logger.Printf(`AddField "%s" for group %s`, name, r.dir)
	if _, ok := r.fields[name]; ok {
		return fmt.Errorf("field \"%s\" has already existed", name)
	}
	r.fields[name] = true
	for _, record := range r.records {
		if err := record.addField(name); err != nil {
			return err
		}
	}
	return nil
}

func (r *ProbDtgp) RemoveField(name string) error {
	logger.Printf(`RemoveField "%s" for group %s`, name, r.dir)
	if _, ok := r.fields[name]; !ok {
		return fmt.Errorf("field \"%s\" not exist", name)
	}
	delete(r.fields, name)
	for _, record := range r.records {
		if err := record.removeField(name); err != nil {
			return err
		}
	}
	return nil
}

// Append an empty record, which will create am empty txt file for each field
func (r *ProbDtgp) NewRecord() (*Record, error) {
	logger.Printf(`NewRecord for group %s`, r.dir)
	record := Record{Map: map[string]string{}, id: len(r.records), dtgp: r}
	r.records = append(r.records, &record)
	for k := range r.fields {
		suf := "txt"
		record.Map[k] = suf // default suffix
		file, err := os.Create(r.rcpath(record.id, k, suf))
		if err != nil {
			return nil, err
		}
		file.Close()
	}
	return &record, nil
}

func (r *ProbDtgp) RemoveRecord(id int) error {
	logger.Printf(`RemoveRecord %d for group %s`, id, r.dir)
	if id < 0 || id > len(r.records) {
		return fmt.Errorf("invalid id")
	}
	r.records[id].Clear()
	r.records = append(r.records[:id], r.records[id+1:]...)
	return nil
}

func LoadDtgp(dirdtgp string) (*ProbDtgp, error) {
	var group ProbDtgp = ProbDtgp{
		dir:     dirdtgp,
		fields:  map[string]bool{},
		records: []*Record{},
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
		group.records = append(group.records, &Record{Map: *record, id: len(group.records), dtgp: &group})
	}
	group.fields = fields

	return &group, nil
}

type Record struct {
	Map  map[string]string
	id   int
	dtgp *ProbDtgp
}

func (r *Record) path(field string) string {
	if _, ok := r.Map[field]; !ok {
		panic("field not exist")
	}
	return r.dtgp.rcpath(r.id, field, r.Map[field])
}
func (r *Record) Clear() {
	for k := range r.Map {
		os.Remove(r.path(k))
	}
}

func (r *Record) addField(name string) error {
	suf := "txt" // default suffix
	r.Map[name] = suf
	file, err := os.Create(r.path(name))
	if err != nil {
		return err
	}
	file.Close()
	return nil
}

func (r *Record) removeField(name string) error {
	os.Remove(r.path(name))
	delete(r.Map, name)
	return nil
}

// Change a field's value to content in filepath, which also update the suffix.
func (r *Record) AlterValue(field string, filepath string) error {
	logger.Printf(`AlterValue id=%d field=%s filepath=%s for group %s`, r.id, field, filepath, r.dtgp.dir)
	if _, ok := r.Map[field]; !ok {
		return fmt.Errorf("field not found")
	}
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

	os.Remove(r.path(field))
	r.Map[field] = ext

	file, err := os.Create(r.path(field))
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

// Get the corresponding record via id, changing suffix to file's path. Return
// nil for invalid id.
func (r *Record) PathMap() *map[string]string {
	a := map[string]string{}
	for field := range r.Map {
		a[field] = r.path(field)
	}
	return &a
}
