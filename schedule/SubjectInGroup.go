package schedule

func (gen *Generator) SubjectInGroup(subject string, group string) bool {
	if !gen.InGroups(group) {
		return false
	}
	if _, ok := gen.Groups[group].Subjects[subject]; ok {
		return true
	}
	return false
}
