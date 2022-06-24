package problem

// map[field]path path 是以 dir 为根目录的相对路径
type record map[string]string

type table struct {
	Field  map[string]bool
	Record []record
}

func (r *table) Fields() *tableFieldCtrl {
	return (*tableFieldCtrl)(r)
}
func (r *table) Records() *tableRecordCtrl {
	return (*tableRecordCtrl)(r)
}

func newTable() table {
	return table{
		Field:  map[string]bool{},
		Record: []record{},
	}
}

type tableFieldCtrl table

// do nothing if name exist
func (r *tableFieldCtrl) Add(name string) {
	if _, ok := r.Field[name]; ok {
		return
	}
	r.Field[name] = true
}

// do nothing if not exist
func (r *tableFieldCtrl) Delete(name string) {
	if _, ok := r.Field[name]; !ok {
		return
	}
	delete(r.Field, name)
	for _, r2 := range r.Record {
		delete(r2, name)
	}
}

type tableRecordCtrl table

func (r *tableRecordCtrl) New() record {
	rcd := record{}
	r.Record = append(r.Record, rcd)
	return rcd
}

func (r *tableRecordCtrl) Delete(id int) {
	r.Record = append(r.Record[:id], r.Record[id+1:]...)
}
