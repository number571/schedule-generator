package schedule

func (gen *Generator) tryGenerate(subgroup SubgroupType, group *Group, subject *Subject, schedule *Schedule) bool {
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
				if 	gen.cellIsReserved(A, schedule, i) && !gen.cellIsReserved(B, schedule, i) && 
					gen.cellIsReserved(B, schedule, i+1) && !gen.cellIsReserved(A, schedule, i+1) ||
					gen.cellIsReserved(B, schedule, i) && !gen.cellIsReserved(A, schedule, i) && 
					gen.cellIsReserved(A, schedule, i+1) && !gen.cellIsReserved(B, schedule, i+1) {
						return true
					}
			}

			// "Подтягивать" полные пары к уже существующим [перед].
			for i := uint8(0); i < gen.NumTables-1; i++ {
				if 	(gen.cellIsReserved(A, schedule, i+1) || gen.cellIsReserved(B, schedule, i+1)) &&
					(!gen.cellIsReserved(A, schedule, i) && !gen.cellIsReserved(B, schedule, i)) {
						lesson = i
						break
					}
			}
		default:
			// "Подтягивать" неполные пары к уже существующим [перед].
			for i := uint8(0); i < gen.NumTables-1; i++ {
				if 	(gen.cellIsReserved(ALL, schedule, i+1) || gen.cellIsReserved(subgroup, schedule, i+1)) &&
					!gen.cellIsReserved(subgroup, schedule, i) {
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
					if 	(gen.cellIsReserved(A, schedule, i) || gen.cellIsReserved(B, schedule, i)) &&
						(!gen.cellIsReserved(A, schedule, i+1) && !gen.cellIsReserved(B, schedule, i+1)) {
							lesson = i+1
							break
						}
				}
			default:
				// "Подтягивать" неполные пары к уже существующим [после].
				for i := uint8(0); i < gen.NumTables-1; i++ {
					if 	(gen.cellIsReserved(ALL, schedule, i) || gen.cellIsReserved(subgroup, schedule, i)) &&
						!gen.cellIsReserved(subgroup, schedule, i+1) {
							lesson = i+1
							break
						}
				}
			}
		}

		cabinet := ""
		if 	gen.cellIsReserved(subgroup, schedule, lesson) || 
			gen.teacherIsReserved(subject.Teacher, lesson) || 
			gen.cabinetIsReserved(subject.Teacher, lesson, &cabinet) {
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
				if 	gen.cellIsReserved(A, schedule, i) && !gen.cellIsReserved(B, schedule, i) && gen.cellIsReserved(ALL, schedule, i+1) ||
					gen.cellIsReserved(B, schedule, i) && !gen.cellIsReserved(A, schedule, i) && gen.cellIsReserved(ALL, schedule, i+1) {
						lesson = savedLesson
						continue nextLesson
					}
			}
		default:
			// Если у одной подгруппы уже имеется пара, а у второй стоит пара
			// в это же время, тогда пропустить проверку пустых окон.
			if 	gen.cellIsReserved(A, schedule, lesson) && A != subgroup || 
				gen.cellIsReserved(B, schedule, lesson) && B != subgroup {
					goto passcheck2
			}
			// Если стоит полная пара, а за ней идёт подгруппа неравная проверяемой, тогда
			// прекратить ставить пары у проверяемой подгруппы.
			fullLessons := false
			for i := uint8(0); i < lesson; i++ {
				if gen.cellIsReserved(ALL, schedule, i) {
					fullLessons = true
					continue
				}
				if 	(fullLessons && gen.cellIsReserved(A, schedule, i) && A != subgroup) ||
					(fullLessons && gen.cellIsReserved(B, schedule, i) && B != subgroup) {
						return true
				}
			}
		}

passcheck2:
		// [ III ] Третья проверка.
		if lesson > 1 {
			// Если нет возможности добавить новые пары без создания окон, тогда не ставить пары.
			for i := uint8(0); i < lesson-1; i++ {
				if 	gen.cellIsReserved(subgroup, schedule, i) && !gen.cellIsReserved(subgroup, schedule, lesson-1) {
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

func (gen *Generator) cellIsReserved(subgroup SubgroupType, schedule *Schedule, lesson uint8) bool {
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

func (gen *Generator) cabinetIsReserved(teacher string, lesson uint8, cabinet *string) bool {
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

func (gen *Generator) teacherIsReserved(teacher string, lesson uint8) bool {
	gen.teacherToReserved(teacher)
	if value, ok := gen.Reserved.Teachers[teacher]; ok {
		return value[lesson] == true
	}
	return false
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
