package schedule

func (gen *Generator) NotHaveLessons(subgroup SubgroupType, subject *Subject, semester uint8) bool {
	switch subgroup {
	case A:
		if subject.Subgroup.A[semester].WeekLessons == 0 {
			return true
		}
	case B:
		if subject.Subgroup.B[semester].WeekLessons == 0 {
			return true
		}
	case ALL:
		if subject.Subgroup.A[semester].WeekLessons == 0 && subject.Subgroup.B[semester].WeekLessons == 0 {
			return true
		}
	}
	return false
}
