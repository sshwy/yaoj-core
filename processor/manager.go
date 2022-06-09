package processor

var processors map[string]Processor = make(map[string]Processor)

func Get(name string) Processor {
	return processors[name]
}

// register a processor to system
func Register(name string, proc Processor) {
	processors[name] = proc
}

func init() {
	Register("checker:hcmp", CheckerHcmp{})
}
