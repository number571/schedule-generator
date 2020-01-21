package schedule

func (gen *Generator) InTeachers(teacher string) bool {
	if _, ok := gen.Teachers[teacher]; ok {
		return true
	}
	return false
}
