package schedule

import (
	
)

func NewGenerator(data *GenData) *Generator {
	if data.Semester > 1 {
		panic("semester /= 0 and /= 1")
	}
	return &Generator{
		Day: data.Day,
		Semester: data.Semester,
		NumTables: data.NumTables,
		Groups: data.Groups,
		Reserved: Reserved{
			Teachers: make(map[string][]bool),
			Cabinets: make(map[string][]bool),
		},
	}
}

func (gen *Generator) Generate() []Schedule {
	var list []Schedule
	for grname, group := range gen.Groups {
		var schedule = gen.newSchedule(grname)
		nextsub: for sbname, subject := range group.Subjects {
			for couple := uint8(0); couple < gen.NumTables; couple++ {
				if (gen.Day != WEDNESDAY && gen.Day != SATURDAY) && couple == 6 {
					break nextsub
				}
				if couple > 1 && isReserved(schedule, couple-2) && !isReserved(schedule, couple-1) {
					break nextsub
				}
				if notHaveHours(subject, gen.Semester) {
					break
				}
				if isReserved(schedule, couple) || gen.teacherIsReserved(subject.Teacher, couple) {
					continue
				}
				cabinet := ""
				if gen.cabinetIsReserved(subject.Teacher, couple, &cabinet) {
					continue
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
	gen.Day = (gen.Day + 1) % 7
	return list
}

func (gen *Generator) newSchedule(group string) Schedule {
	return Schedule{
		Day: gen.Day,
		Group: group,
		Table: make([]Row, gen.NumTables),
	}
}

func (gen *Generator) teacherToReserved(teacher string) {
	if _, ok := gen.Reserved.Teachers[teacher]; ok {
		return
	}
	gen.Reserved.Teachers[teacher] = make([]bool, gen.NumTables)
}

func (gen *Generator) cabinetToReserved(cabnum string) {
	if _, ok := gen.Reserved.Cabinets[cabnum]; ok {
		return
	}
	gen.Reserved.Cabinets[cabnum] = make([]bool, gen.NumTables)
}

func (gen *Generator) teacherIsReserved(teacher string, couple uint8) bool {
	gen.teacherToReserved(teacher)
	if value, ok := gen.Reserved.Teachers[teacher]; ok {
		return value[couple] == true
	}
	return false
}

func (gen *Generator) cabinetIsReserved(teacher string, couple uint8, cabinet *string) bool {
	var result = true
	for cabnum := range Teachers[teacher].Cabinets {
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

func isReserved(schedule Schedule, couple uint8) bool {
	if schedule.Table[couple].Subject == "" {
		return false
	}
	return true
}

func notHaveHours(subject LocalSubj, semester uint8) bool {
	if subject.Hours.Semester[semester].WeekHours == 0 {
		return true
	}
	return false
}
