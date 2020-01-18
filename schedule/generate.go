package schedule

import (
	// "fmt"
	"errors"
)

func NewGenerator(data *Generator) *Generator {
	if data.Semester != 1 && data.Semester != 2 {
		panic("semester /= 1 and /= 2")
	}
	data.Semester -= 1
	return &Generator{
		Day: data.Day,
		Semester: data.Semester,
		NumTables: data.NumTables,
		Groups: data.Groups,
		Teachers: data.Teachers,
		Blocked: make(map[string]bool),
		Reserved: Reserved{
			Teachers: make(map[string][]bool),
			Cabinets: make(map[string][]bool),
		},
	}
}

func (gen *Generator) NewSchedule(group string) *Schedule {
	return &Schedule{
		Day: gen.Day,
		Group: group,
		Table: make([]Row, gen.NumTables),
	}
}

func (gen *Generator) Generate() []*Schedule {
	var list []*Schedule
	for _, group := range gen.Groups {
		var schedule = gen.NewSchedule(group.Name)
		if gen.Day == SUNDAY {
			list = append(list, schedule)
			continue
		}
		nextsub: for _, subject := range group.Subjects {
			if !subject.IsSplited {
				// fmt.Println("not splited")
				if gen.generateSubgroup(ALL, group, subject, schedule) {
					break nextsub
				}
			} else {
				// fmt.Println("splited")
				if gen.generateSubgroup(A, group, subject, schedule) {
					break nextsub
				}
				if gen.generateSubgroup(B, group, subject, schedule) {
					break nextsub
				}
			}
		}
		list = append(list, schedule)
	}
	gen.Reserved.Teachers = make(map[string][]bool)
	gen.Reserved.Cabinets = make(map[string][]bool)
	gen.Day = (gen.Day + 1) % 7
	return list
}

func (gen *Generator) generateSubgroup(subgroup SubgroupType, group *Group, subject *Subject, schedule *Schedule) bool {
	tryBreak := false

	for couple := uint8(0); couple < gen.NumTables; couple++ {
		if gen.InBlocked(subject.Teacher) || gen.NotHaveHours(subgroup, subject, gen.Semester){
			break
		}

		saveCouple := couple

		// Without middle spaces.
		if subgroup == ALL {
			for i := uint8(1); i < gen.NumTables; i++ {
				if 	(gen.CellIsReserved(A, schedule, i) || gen.CellIsReserved(B, schedule, i)) &&
					!gen.CellIsReserved(ALL, schedule, i-1) {
						couple = i-1
						tryBreak = true
						break
				}
			}
		}

		cabinet := ""
		if (gen.CellIsReserved(ALL, schedule, couple) || 
			gen.TeacherIsReserved(subject.Teacher, couple) || 
			gen.CabinetIsReserved(subject.Teacher, couple, &cabinet)) && tryBreak {
				return true
		}

		if 	gen.CellIsReserved(ALL, schedule, couple) || 
			gen.TeacherIsReserved(subject.Teacher, couple) || 
			gen.CabinetIsReserved(subject.Teacher, couple, &cabinet){
			continue
		}

		// Full day = max 6 couples.
		if (gen.Day != WEDNESDAY && gen.Day != SATURDAY) && couple == 6 {
			if subgroup == ALL {
				return true
			}
			break
		}

		// Without middle spaces.
		if couple > 1 {
			switch subgroup {
			case ALL:
				for i := uint8(0); i < couple; i++ {
					if 	gen.CellIsReserved(A, schedule, i) && gen.CellIsReserved(B, schedule, couple-1) ||
						gen.CellIsReserved(B, schedule, i) && gen.CellIsReserved(A, schedule, couple-1) {
							return true
					}
				}
			default:
				for i := uint8(0); i < couple; i++ {
					if 	gen.CellIsReserved(subgroup, schedule, i) && !gen.CellIsReserved(subgroup, schedule, couple-1)  {
							return true
					}
				}
			}
		}

		gen.Reserved.Teachers[subject.Teacher][couple] = true
		gen.Reserved.Cabinets[cabinet][couple] = true

		switch subgroup{
		case A: gen.Groups[group.Name].Subjects[subject.Name].Subgroup.A[gen.Semester].WeekHours -= 2
		case B: gen.Groups[group.Name].Subjects[subject.Name].Subgroup.B[gen.Semester].WeekHours -= 2
		case ALL:
			gen.Groups[group.Name].Subjects[subject.Name].Subgroup.A[gen.Semester].WeekHours -= 2
			gen.Groups[group.Name].Subjects[subject.Name].Subgroup.B[gen.Semester].WeekHours -= 2
		}

		if subgroup == ALL {
			schedule.Table[couple].Teacher = [2]string{
				subject.Teacher,
				subject.Teacher,
			}
			schedule.Table[couple].Subject = [2]string{
				subject.Name,
				subject.Name,
			}
			schedule.Table[couple].Cabinet = [2]string{
				cabinet,
				cabinet,
			}
			if saveCouple != couple {
				couple = saveCouple
			}
			continue
		}

		schedule.Table[couple].Teacher[subgroup] = subject.Teacher
		schedule.Table[couple].Subject[subgroup] = subject.Name
		schedule.Table[couple].Cabinet[subgroup] = cabinet
	}
	return false
}

func (gen *Generator) SubjectInGroup(subject string, group string) bool {
	if !gen.InGroups(group) {
		return false
	}
	if _, ok := gen.Groups[group].Subjects[subject]; ok {
		return true
	}
	return false
}

func (gen *Generator) InGroups(group string) bool {
	if _, ok := gen.Groups[group]; ok {
		return true
	}
	return false
}

func (gen *Generator) InTeachers(teacher string) bool {
	if _, ok := gen.Teachers[teacher]; ok {
		return true
	}
	return false
}

func (gen *Generator) InBlocked(teacher string) bool {
	if _, ok := gen.Blocked[teacher]; ok {
		return true
	}
	return false
}

func (gen *Generator) BlockTeacher(teacher string) error {
	if !gen.InTeachers(teacher) {
		return errors.New("teacher undefined")
	}
	gen.Blocked[teacher] = true
	return nil
}

func (gen *Generator) UnblockTeacher(teacher string) error {
	if !gen.InBlocked(teacher) {
		return errors.New("teacher undefined")
	}
	delete(gen.Blocked, teacher)
	return nil
}

func (gen *Generator) TeacherIsReserved(teacher string, couple uint8) bool {
	gen.teacherToReserved(teacher)
	if value, ok := gen.Reserved.Teachers[teacher]; ok {
		return value[couple] == true
	}
	return false
}

func (gen *Generator) CabinetIsReserved(teacher string, couple uint8, cabinet *string) bool {
	var result = true
	for _, cabnum := range gen.Teachers[teacher].Cabinets {
		gen.cabinetToReserved(cabnum)
		if _, ok := gen.Reserved.Cabinets[cabnum]; ok {
			if gen.Reserved.Cabinets[cabnum][couple] == false {
				*cabinet = cabnum
				return false
			}
		}
	}
	return result
}

func (gen *Generator) CellIsReserved(subgroup SubgroupType, schedule *Schedule, couple uint8) bool {
	switch subgroup {
	case A: 
		if schedule.Table[couple].Subject[A] == "" {
			return false
		}
	case B:
		if schedule.Table[couple].Subject[B] == "" {
			return false
		}
	case ALL:
		if schedule.Table[couple].Subject[A] == "" && schedule.Table[couple].Subject[B] == "" {
			return false
		}
	}
	return true
}

func (gen *Generator) NotHaveHours(subgroup SubgroupType, subject *Subject, semester uint8) bool {
	switch subgroup {
	case A:
		if subject.Subgroup.A[semester].WeekHours == 0 {
			return true
		}
	case B:
		if subject.Subgroup.B[semester].WeekHours == 0 {
			return true
		}
	case ALL:
		if subject.Subgroup.A[semester].WeekHours == 0 && subject.Subgroup.B[semester].WeekHours == 0 {
			return true
		}
	}
	return false
}
