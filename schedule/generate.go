package schedule

import (
	"fmt"
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
				if group.Name == "102" {
					fmt.Println("not splited")
				}
				
				if gen.generateSubgroup(ALL, group, subject, schedule) {
					break nextsub
				}
			} else {
				if group.Name == "102" {
					fmt.Println("splited")
				}

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
	nextLesson: for lesson := uint8(0); lesson < gen.NumTables; lesson++ {
		if gen.InBlocked(subject.Teacher) || gen.NotHaveLessons(subgroup, subject, gen.Semester){
			break
		}

		isAfter := false
		savedLesson := lesson

		// [ I ] Первая проверка.
		switch subgroup {
		case ALL:
			// Если две подгруппы стоят друг за другом, тогда исключить возможность добавления полной пары.
			for i := uint8(0); i < gen.NumTables-1; i++ {
				if 	gen.CellIsReserved(A, schedule, i) && !gen.CellIsReserved(B, schedule, i) && gen.CellIsReserved(B, schedule, i+1) && !gen.CellIsReserved(A, schedule, i+1) ||
					gen.CellIsReserved(B, schedule, i) && !gen.CellIsReserved(A, schedule, i) && gen.CellIsReserved(A, schedule, i+1) && !gen.CellIsReserved(B, schedule, i+1) {
						return true
					}
			}
			// "Подтягивать" полные пары к уже существующим [перед].
			for i := uint8(0); i < gen.NumTables-1; i++ {
				if 	(gen.CellIsReserved(A, schedule, i+1) || gen.CellIsReserved(B, schedule, i+1)) &&
					!gen.CellIsReserved(ALL, schedule, i) {
						lesson = i
						break
					}
			}
		default:
			// "Подтягивать" неполные пары к уже существующим полным [перед].
			for i := uint8(0); i < gen.NumTables-1; i++ {
				if 	(gen.CellIsReserved(ALL, schedule, i+1) || gen.CellIsReserved(subgroup, schedule, i+1)) &&
					!gen.CellIsReserved(subgroup, schedule, i) {
						lesson = i
						break
					}
			}
		}

tryAfter:
		if isAfter {
			switch subgroup {
			case ALL:
				// "Подтягивать" полные пары к уже существующим [после].
				for i := uint8(0); i < gen.NumTables-1; i++ {
					if 	(gen.CellIsReserved(A, schedule, i) || gen.CellIsReserved(B, schedule, i)) &&
						!gen.CellIsReserved(ALL, schedule, i+1) {
							lesson = i+1
							break
						}
				}
			default:
				// "Подтягивать" неполные пары к уже существующим [после].
				for i := uint8(0); i < gen.NumTables-1; i++ {
					if 	(gen.CellIsReserved(ALL, schedule, i) || gen.CellIsReserved(subgroup, schedule, i)) &&
						!gen.CellIsReserved(subgroup, schedule, i+1) {
							lesson = i+1
							break
						}
				}
			}
			
		}

		cabinet := ""
		if 	gen.CellIsReserved(subgroup, schedule, lesson) || 
			gen.TeacherIsReserved(subject.Teacher, lesson) || 
			gen.CabinetIsReserved(subject.Teacher, lesson, &cabinet) {
				if isAfter {
					break nextLesson
				}
				if lesson != savedLesson {
					isAfter = true
					goto tryAfter
				}
				continue nextLesson
		}

		// Full day = max 6 couples.
		if (gen.Day != WEDNESDAY && gen.Day != SATURDAY) && lesson == 6 {
			break nextLesson
		}

		// [ II ] Вторая проверка.
		switch subgroup {
		case ALL:
			// Если уже существует полная пара, которая стоит за парами с подгруппами, тогда
			// перейти на новую ячейку расписания группы.
			for i := lesson; i < gen.NumTables-1; i++ {
				if 	gen.CellIsReserved(A, schedule, i) && !gen.CellIsReserved(B, schedule, i) && gen.CellIsReserved(ALL, schedule, i+1) ||
					gen.CellIsReserved(B, schedule, i) && !gen.CellIsReserved(A, schedule, i) && gen.CellIsReserved(ALL, schedule, i+1) {
						lesson = savedLesson
						continue nextLesson
					}
			}
		default:
			// Если у одной подгруппы уже имеется пара, а у второй стоит пара
			// в это же время, тогда пропустить проверку пустых окон.
			if 	gen.CellIsReserved(A, schedule, lesson) && A != subgroup || 
				gen.CellIsReserved(B, schedule, lesson) && B != subgroup {
					goto passcheck2
			}
			// Если стоит полная пара, а за ней идёт подгруппа неравная проверяемой, тогда
			// прекратить ставить пары у проверяемой подгруппы.
			fullLessons := false
			for i := uint8(0); i < lesson; i++ {
				if gen.CellIsReserved(ALL, schedule, i) {
					fullLessons = true
					continue
				}
				if 	(fullLessons && gen.CellIsReserved(A, schedule, i) && A != subgroup) ||
					(fullLessons && gen.CellIsReserved(B, schedule, i) && B != subgroup) {
						return true
				}
			}
		}

passcheck2:
		// [ III ] Третья проверка.
		if lesson > 1 {
			// Если нет возможности добавить новые пары без создания окон, тогда не ставить пары.
			for i := uint8(0); i < lesson-1; i++ {
				if gen.CellIsReserved(subgroup, schedule, i) && !gen.CellIsReserved(subgroup, schedule, lesson-1) {
					return false
				}
			}
		}

		gen.Reserved.Teachers[subject.Teacher][lesson] = true
		gen.Reserved.Cabinets[cabinet][lesson] = true

		switch subgroup{
		case A: gen.Groups[group.Name].Subjects[subject.Name].Subgroup.A[gen.Semester].WeekLessons -= 1
		case B: gen.Groups[group.Name].Subjects[subject.Name].Subgroup.B[gen.Semester].WeekLessons -= 1
		case ALL:
			gen.Groups[group.Name].Subjects[subject.Name].Subgroup.A[gen.Semester].WeekLessons -= 1
			gen.Groups[group.Name].Subjects[subject.Name].Subgroup.B[gen.Semester].WeekLessons -= 1
		}

		if subgroup == ALL {
			schedule.Table[lesson].Teacher = [2]string{
				subject.Teacher,
				subject.Teacher,
			}
			schedule.Table[lesson].Subject = [2]string{
				subject.Name,
				subject.Name,
			}
			schedule.Table[lesson].Cabinet = [2]string{
				cabinet,
				cabinet,
			}
			lesson = savedLesson
			continue nextLesson
		}

		schedule.Table[lesson].Teacher[subgroup] = subject.Teacher
		schedule.Table[lesson].Subject[subgroup] = subject.Name
		schedule.Table[lesson].Cabinet[subgroup] = cabinet
		lesson = savedLesson
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

func (gen *Generator) TeacherIsReserved(teacher string, lesson uint8) bool {
	gen.teacherToReserved(teacher)
	if value, ok := gen.Reserved.Teachers[teacher]; ok {
		return value[lesson] == true
	}
	return false
}

func (gen *Generator) CabinetIsReserved(teacher string, lesson uint8, cabinet *string) bool {
	var result = true
	for _, cabnum := range gen.Teachers[teacher].Cabinets {
		gen.cabinetToReserved(cabnum)
		if _, ok := gen.Reserved.Cabinets[cabnum]; ok {
			if gen.Reserved.Cabinets[cabnum][lesson] == false {
				*cabinet = cabnum
				return false
			}
		}
	}
	return result
}

func (gen *Generator) CellIsReserved(subgroup SubgroupType, schedule *Schedule, lesson uint8) bool {
	switch subgroup {
	case A: 
		if schedule.Table[lesson].Subject[A] != "" {
			return true
		}
	case B:
		if schedule.Table[lesson].Subject[B] != "" {
			return true
		}
	case ALL:
		if schedule.Table[lesson].Subject[A] != "" && schedule.Table[lesson].Subject[B] != "" {
			return true
		}
	}
	return false
}

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
