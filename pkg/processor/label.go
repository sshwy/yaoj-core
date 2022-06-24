package processor

var inLabel, ouLabel map[string][]string = map[string][]string{}, map[string][]string{}

func InputLabel(name string) []string {
	return inLabel[name]
}
func OutputLabel(name string) []string {
	return ouLabel[name]
}
