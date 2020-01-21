package schedule

func (gen *Generator) InGroups(group string) bool {
	if _, ok := gen.Groups[group]; ok {
		return true
	}
	return false
}
