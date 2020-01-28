package schedule

import (
    "os"
    "sort"
    "time"
    "errors"
    "math/rand"
    "encoding/json"
)

func init() {
    rand.Seed(time.Now().UnixNano())
}

func (gen *Generator) tryGenerate(subgroup SubgroupType, group *Group, subject *Subject, schedule *Schedule, subjtype SubjectType) {
    nextLesson: for lesson := uint(0); lesson < gen.NumTables; lesson++ {

        // Если это полная группа и у неё не осталось теоретических занятий, тогда пропустить этот предмет. 
        if subjtype == THEORETICAL && !gen.haveTheoreticalLessons(subject) {
            break nextLesson
        }

        // Если это подгруппа и у неё не осталось практических занятий, тогда пропустить этот предмет.
        if subjtype == PRACTICAL && !gen.havePracticalLessons(subgroup, subject) {
            break nextLesson
        }

        // Если учитель заблокирован (не может присутствовать на занятиях) или
        // лимит пар на текущую неделю для предмета завершён, тогда пропустить этот предмет.
        if gen.inBlocked(subject.Teacher) || gen.notHaveWeekLessons(subgroup, subject) {
            break nextLesson
        }

        isAfter := false
        savedLesson := lesson

        // [ I ] Первая проверка.
        switch subgroup {
        case ALL:
            // Если две подгруппы стоят друг за другом, тогда исключить возможность добавления полной пары.
            for i := uint(0); i < gen.NumTables-1; i++ {
                if  gen.cellIsReserved(A, schedule, i) && !gen.cellIsReserved(B, schedule, i) && 
                    gen.cellIsReserved(B, schedule, i+1) && !gen.cellIsReserved(A, schedule, i+1) ||
                    gen.cellIsReserved(B, schedule, i) && !gen.cellIsReserved(A, schedule, i) && 
                    gen.cellIsReserved(A, schedule, i+1) && !gen.cellIsReserved(B, schedule, i+1) {
                        break nextLesson
                    }
            }

            // Если между двумя разными подгруппами находятся окна, тогда посчитать сколько пустых окон.
            // Если их количество = 1, тогда попытаться подставить полную пару под это окно. 
            // Если не получается, тогда не ставить полную пару.
            cellSubgroupReserved := false
            indexNotReserved := uint(0)
            cellIsNotReserved := 0
            for i := uint(0); i < gen.NumTables-1; i++ {
                if  gen.cellIsReserved(A, schedule, i) && !gen.cellIsReserved(B, schedule, i) ||
                    gen.cellIsReserved(B, schedule, i) && !gen.cellIsReserved(A, schedule, i) {
                        if cellIsNotReserved != 0 {
                            if cellIsNotReserved == 1 {
                                lesson = indexNotReserved
                                goto tryAfter
                            }
                            break nextLesson
                        }
                        cellSubgroupReserved = true
                    }
                if !gen.cellIsReserved(A, schedule, i) && !gen.cellIsReserved(B, schedule, i) && cellSubgroupReserved {
                    if cellIsNotReserved == 0 {
                        indexNotReserved = i
                    }
                    cellIsNotReserved += 1
                }
            }

            // "Подтягивать" полные пары к уже существующим [перед].
            for i := uint(0); i < gen.NumTables-1; i++ {
                if  (gen.cellIsReserved(A, schedule, i+1) || gen.cellIsReserved(B, schedule, i+1)) &&
                    (!gen.cellIsReserved(A, schedule, i) && !gen.cellIsReserved(B, schedule, i)) {
                        lesson = i
                        break
                    }
            }
        default:
            // "Подтягивать" неполные пары к уже существующим [перед].
            for i := uint(0); i < gen.NumTables-1; i++ {
                if  !gen.cellIsReserved(subgroup, schedule, i) && gen.cellIsReserved(subgroup, schedule, i+1) {
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
                for i := uint(0); i < gen.NumTables-1; i++ {
                    if  (gen.cellIsReserved(A, schedule, i) || gen.cellIsReserved(B, schedule, i)) &&
                        (!gen.cellIsReserved(A, schedule, i+1) && !gen.cellIsReserved(B, schedule, i+1)) {
                            lesson = i+1
                            break
                        }
                }
            default:
                // "Подтягивать" неполные пары к уже существующим [после].
                for i := uint(0); i < gen.NumTables-1; i++ {
                    if  (gen.cellIsReserved(ALL, schedule, i) || gen.cellIsReserved(subgroup, schedule, i)) &&
                        !gen.cellIsReserved(subgroup, schedule, i+1) {
                            lesson = i+1
                            break
                        }
                }
            }
        }

        var (
            cabinet = ""
            cabinet2 = ""
        )
        // Если ячейка уже зарезервирована или учитель занят, или кабинет зарезервирован, или
        // если это полная пара и ячейка занята либо первой, либо второй подгруппой, или
        // если это двойная пара и кабинеты второго учителя зарезервированы, или 
        // если это двойная пара и второй учитель занят: попытаться сдвинуть пары к уже существующим.
        // Если не получается, тогда не ставить этот урок.
        if  gen.cellIsReserved(subgroup, schedule, lesson) || 
            gen.teacherIsReserved(subject.Teacher, lesson) || 
            gen.cabinetIsReserved(subject.Teacher, lesson, &cabinet) || 
            (subgroup == ALL && (gen.cellIsReserved(A, schedule, lesson) || gen.cellIsReserved(B, schedule, lesson))) ||
            (gen.withSubgroups(group.Name) && gen.isDoubleLesson(group.Name, subject.Name) && gen.cabinetIsReserved(subject.Teacher2, lesson, &cabinet2)) ||
            (gen.withSubgroups(group.Name) && gen.isDoubleLesson(group.Name, subject.Name) && gen.teacherIsReserved(subject.Teacher2, lesson)) {
                if isAfter {
                    break nextLesson
                }
                if lesson != savedLesson {
                    isAfter = true
                    goto tryAfter
                }
                continue nextLesson
        }

        // Полный день - максимум 6 пар.
        // lesson начинается с нуля!
        if (gen.Day != WEDNESDAY && gen.Day != SATURDAY) && lesson >= 6 {
            break nextLesson
        }

        // В среду и субботу - 5-6 пары должны быть свободны.
        // lesson начинается с нуля!
        if (gen.Day == WEDNESDAY || gen.Day == SATURDAY) && (lesson == 4 || lesson == 5) {
            lesson = savedLesson // Возобновить сохранённое занятие.
            continue nextLesson // Перейти к следующей ячейке.
        }

        // [ II ] Вторая проверка.
        switch subgroup {
        case ALL:
            // Если уже существует полная пара, которая стоит за парами с подгруппами, тогда
            // перейти на новую ячейку расписания группы.
            for i := lesson; i < gen.NumTables-1; i++ {
                if  gen.cellIsReserved(A, schedule, i) && !gen.cellIsReserved(B, schedule, i) && gen.cellIsReserved(ALL, schedule, i+1) ||
                    gen.cellIsReserved(B, schedule, i) && !gen.cellIsReserved(A, schedule, i) && gen.cellIsReserved(ALL, schedule, i+1) {
                        lesson = savedLesson
                        continue nextLesson
                    }
            }
        default:
            // Если у одной подгруппы уже имеется пара, а у второй стоит пара
            // в это же время, тогда пропустить проверку пустых окон.
            if  gen.cellIsReserved(A, schedule, lesson) && A != subgroup || 
                gen.cellIsReserved(B, schedule, lesson) && B != subgroup {
                    goto passcheck2
            }
            // Если стоит полная пара, а за ней идёт подгруппа неравная проверяемой, тогда
            // прекратить ставить пары у проверяемой подгруппы.
            fullLessons := false
            for i := uint(0); i < lesson; i++ {
                if gen.cellIsReserved(ALL, schedule, i) {
                    fullLessons = true
                    continue
                }
                if  (fullLessons && gen.cellIsReserved(A, schedule, i) && A != subgroup) ||
                    (fullLessons && gen.cellIsReserved(B, schedule, i) && B != subgroup) {
                        break nextLesson
                }
            }
        }

passcheck2:
        // [ III ] Третья проверка.
        // Если нет возможности добавить новые пары без создания окон, тогда не ставить пары.
        if lesson > 1 {
            switch subgroup {
            case ALL:
                for i := uint(0); i < lesson-1; i++ {
                    if  (gen.cellIsReserved(A, schedule, i) && !gen.cellIsReserved(A, schedule, lesson-1)) || 
                        (gen.cellIsReserved(B, schedule, i) && !gen.cellIsReserved(B, schedule, lesson-1)) {
                            break nextLesson
                    }
                }
            default:
                for i := uint(0); i < lesson-1; i++ {
                    if  gen.cellIsReserved(subgroup, schedule, i) && !gen.cellIsReserved(subgroup, schedule, lesson-1) {
                        break nextLesson
                    }
                }
            }
        }

        gen.Reserved.Teachers[subject.Teacher][lesson] = true
        gen.Reserved.Cabinets[cabinet][lesson] = true

        switch subgroup {
        case A:
            gen.Groups[group.Name].Subjects[subject.Name].WeekLessons.A -= 1
            gen.Groups[group.Name].Subjects[subject.Name].Practice.A -= 1
        case B:
            gen.Groups[group.Name].Subjects[subject.Name].WeekLessons.B -= 1
            gen.Groups[group.Name].Subjects[subject.Name].Practice.B -= 1
        case ALL:
            gen.Groups[group.Name].Subjects[subject.Name].WeekLessons.A -= 1
            gen.Groups[group.Name].Subjects[subject.Name].WeekLessons.B -= 1
            switch subjtype {
            case THEORETICAL:
                gen.Groups[group.Name].Subjects[subject.Name].Theory -= 1
            case PRACTICAL:
                gen.Groups[group.Name].Subjects[subject.Name].Practice.A -= 1
                gen.Groups[group.Name].Subjects[subject.Name].Practice.B -= 1
            }
        }

        if subgroup == ALL {
            schedule.Table[lesson].Teacher = [ALL]string{
                subject.Teacher,
                subject.Teacher,
            }
            schedule.Table[lesson].Subject = [ALL]string{
                subject.Name,
                subject.Name,
            }
            schedule.Table[lesson].Cabinet = [ALL]string{
                cabinet,
                cabinet,
            }
            // Если это двойная пара и группа делимая, тогда поставить пару с разными преподавателями.
            if gen.isDoubleLesson(group.Name, subject.Name) && gen.withSubgroups(group.Name) {
                gen.Reserved.Teachers[subject.Teacher2][lesson] = true
                schedule.Table[lesson].Teacher[B] = subject.Teacher2
                schedule.Table[lesson].Cabinet[B] = cabinet2
            }
            lesson = savedLesson
            continue nextLesson
        }

        schedule.Table[lesson].Teacher[subgroup] = subject.Teacher
        schedule.Table[lesson].Subject[subgroup] = subject.Name
        schedule.Table[lesson].Cabinet[subgroup] = cabinet
        lesson = savedLesson
    }
}

func (gen *Generator) withSubgroups(group string) bool {
    if gen.Groups[group].Quantity > MAX_COUNT_WITHOUT_SUBGROUPS {
        return true
    }
    return false
}

func (gen *Generator) blockTeacher(teacher string) error {
    if !gen.inTeachers(teacher) {
        return errors.New("teacher undefined")
    }
    gen.Blocked[teacher] = true
    return nil
}

func (gen *Generator) inBlocked(teacher string) bool {
    if _, ok := gen.Blocked[teacher]; !ok {
        return false
    }
    return true
}

func (gen *Generator) inGroups(group string) bool {
    if _, ok := gen.Groups[group]; !ok {
        return false
    }
    return true
}

func (gen *Generator) inTeachers(teacher string) bool {
    if _, ok := gen.Teachers[teacher]; !ok {
        return false
    }
    return true
}

func (gen *Generator) unblockTeacher(teacher string) error {
    if !gen.inBlocked(teacher) {
        return errors.New("teacher undefined")
    }
    delete(gen.Blocked, teacher)
    return nil
}

func packJSON(data interface{}) []byte {
    jsonData, err := json.MarshalIndent(data, "", "\t")
    if err != nil {
        return nil
    }
    return jsonData
}

func writeJSON(filename string, data interface{}) error {
    return writeFile(filename, string(packJSON(data)))
}

func (gen *Generator) subjectInGroup(subject string, group string) bool {
    if !gen.inGroups(group) {
        return false
    }
    if _, ok := gen.Groups[group].Subjects[subject]; !ok {
        return false
    }
    return true
}

func (gen *Generator) isDoubleLesson(group string, subject string) bool {
    if !gen.inGroups(group) {
        return false
    }
    if _, ok := gen.Groups[group].Subjects[subject]; !ok {
        return false
    }
    if gen.Groups[group].Subjects[subject].Teacher2 == "" {
        return false
    }
    return true
}

func shuffle(slice interface{}) interface{}{
    switch slice.(type) {
    case []*Group:
        result := slice.([]*Group)
        for i := len(result)-1; i > 0; i-- {
            j := rand.Intn(i+1)
            result[i], result[j] = result[j], result[i]
        }
        return result
    case []*Subject:
        result := slice.([]*Subject)
        for i := len(result)-1; i > 0; i-- {
            j := rand.Intn(i+1)
            result[i], result[j] = result[j], result[i]
        }
        return result
    }
    return nil
}

func (gen *Generator) cellIsReserved(subgroup SubgroupType, schedule *Schedule, lesson uint) bool {
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

func (gen *Generator) cabinetIsReserved(teacher string, lesson uint, cabinet *string) bool {
    var result = true
    if !gen.inTeachers(teacher) {
        return result
    }
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

func (gen *Generator) teacherIsReserved(teacher string, lesson uint) bool {
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

func (gen *Generator) notHaveWeekLessons(subgroup SubgroupType, subject *Subject) bool {
    switch subgroup {
    case A:
        if subject.WeekLessons.A == 0 {
            return true
        }
    case B:
        if subject.WeekLessons.B == 0 {
            return true
        }
    case ALL:
        if subject.WeekLessons.A == 0 && subject.WeekLessons.B == 0 {
            return true
        }
    }
    return false
}

func (gen *Generator) haveTheoreticalLessons(subject *Subject) bool {
    if subject.Theory == 0 {
        return false
    }
    return true
}

func (gen *Generator) havePracticalLessons(subgroup SubgroupType, subject *Subject) bool {
    switch subgroup {
    case A:
        if subject.Practice.A == 0 {
            return false
        }
    case B:
        if subject.Practice.B == 0 {
            return false
        }
    case ALL:
        if subject.Practice.A == 0 && subject.Practice.B == 0 {
            return false
        }
    }
    return true
}

func getGroups(groups map[string]*Group) []*Group {
    var list []*Group
    for _, group := range groups {
        list = append(list, group)
    }
    return shuffle(list).([]*Group)
}

func getSubjects(subjects map[string]*Subject) []*Subject {
    var list []*Subject
    for _, subject := range subjects {
        list = append(list, subject)
    }
    return shuffle(list).([]*Subject)
}

func sortSchedule(schedule []*Schedule) []*Schedule {
    sort.SliceStable(schedule, func(i, j int) bool {
        return schedule[i].Group < schedule[j].Group
    })
    return schedule
}

func colWidthForCabinets(index int) (int, int, float64) {
    var col = (index+1)*3+1
    return col, col, COL_W_CAB
}

// Returns [min:max] value.
func random(min, max int) int {
    return rand.Intn(max - min + 1) + min
}

func readFile(filename string) string {
    file, err := os.Open(filename)
    if err != nil {
        return ""
    }
    defer file.Close()

    var (
        buffer []byte = make([]byte, BUFFER)
        data string
    )

    for {
        length, err := file.Read(buffer)
        if length == 0 || err != nil { break }
        data += string(buffer[:length])
    }

    return data
}

func writeFile(filename, data string) error {
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    file.WriteString(data)
    file.Close()
    return nil
}
