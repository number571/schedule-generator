package schedule

import (
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

func (gen *Generator) Generate(condition func(*Generator, *Group, *Subject, *Schedule, *uint8) bool) []*Schedule {
	var list []*Schedule
	for grname, group := range gen.Groups {
		var schedule = gen.NewSchedule(grname)
		if gen.Day == SUNDAY {
			list = append(list, schedule)
			continue	
		}
		nextsub: for sbname, subject := range group.Subjects {
			for couple := uint8(0); couple < gen.NumTables; couple++ {
				if gen.InBlocked(subject.Teacher) || gen.NotHaveHours(subject, gen.Semester){
					break
				}

				cabinet := ""
				if 	gen.IsReserved(schedule, couple) || 
					gen.TeacherIsReserved(subject.Teacher, couple) || 
					gen.CabinetIsReserved(subject.Teacher, couple, &cabinet){
					continue
				}

				if !condition(gen, &group, subject, schedule, &couple) {
					break nextsub
				}

				gen.Groups[grname].Subjects[sbname].Hours.Semester[gen.Semester].WeekHours -= 2

				gen.Reserved.Teachers[subject.Teacher][couple] = true
				gen.Reserved.Cabinets[cabinet][couple] = true

				schedule.Table[couple].Teacher = subject.Teacher
				schedule.Table[couple].Subject = sbname
				schedule.Table[couple].Cabinet = cabinet
			}
		}
		list = append(list, schedule)
	}
	gen.Reserved.Teachers = make(map[string][]bool)
	gen.Reserved.Cabinets = make(map[string][]bool)
	gen.Day = (gen.Day + 1) % 7
	return list
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

func (gen *Generator) IsReserved(schedule *Schedule, couple uint8) bool {
	if schedule.Table[couple].Subject == "" {
		return false
	}
	return true
}

func (gen *Generator) NotHaveHours(subject *Subject, semester uint8) bool {
	if subject.Hours.Semester[semester].WeekHours == 0 {
		return true
	}
	return false
}
