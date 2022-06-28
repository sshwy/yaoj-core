package run

type WorkflowCache interface {
	// hash value of node and its output files
	Set(hash sha, outputs []string)
	// nil if no cache
	Get(hash sha) []string
}

type InMemoryCache struct {
	data map[sha][]string
}

var _ WorkflowCache = (*InMemoryCache)(nil)

func (r *InMemoryCache) Set(hash sha, outputs []string) {
	if r.data[hash] != nil {
		panic("multi set")
	}
	// logger.Print("\033[32mSet: \033[0m", hash, outputs)
	r.data[hash] = outputs[:]
}

func (r *InMemoryCache) Get(hash sha) []string {
	return r.data[hash][:]
}

var globalCache = InMemoryCache{
	data: map[sha][]string{},
}
