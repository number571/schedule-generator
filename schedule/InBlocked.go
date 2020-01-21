package schedule

func (gen *Generator) InBlocked(teacher string) bool {
	if _, ok := gen.Blocked[teacher]; ok {
		return true
	}
	return false
}
