package schedule

import (
    "fmt"
    "strconv"
    "encoding/json"
    "github.com/tealeg/xlsx"
)

func NewGenerator(data *Generator) *Generator {
    return &Generator{
        Day: data.Day,
        Debug: data.Debug,
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
        Table: make([]Row, NUM_TABLES),
    }
}

func ReadGroups(filename string) map[string]*Group {
    var (
        groups = make(map[string]*Group)
        groupsList []GroupJSON
    )
    data := readFile(filename)
    err := json.Unmarshal([]byte(data), &groupsList)
    if err != nil {
        return nil
    }
    for _, gr := range groupsList {
        groups[gr.Name] = &Group{
            Name: gr.Name,
            Quantity: gr.Quantity,
        }
        groups[gr.Name].Subjects = make(map[string]*Subject)
        for _, sb := range gr.Subjects {
            if _, ok := groups[gr.Name].Subjects[sb.Name]; ok {
                groups[gr.Name].Subjects[sb.Name].Teacher2 = sb.Teacher
                continue
            }
            groups[gr.Name].Subjects[sb.Name] = &Subject{
                Name: sb.Name,
                Teacher: sb.Teacher,
                IsComputer: sb.IsComputer,
                SaveWeek: sb.Lessons.Week,
                Theory: sb.Lessons.Theory,
                Practice: Subgroup{
                    A: sb.Lessons.Practice,
                    B: sb.Lessons.Practice,
                },
                WeekLessons: Subgroup{
                    A: sb.Lessons.Week,
                    B: sb.Lessons.Week,
                },
            }
        }
    }
    return groups
}

func ReadTeachers(filename string) map[string]*Teacher {
    var (
        teachers = make(map[string]*Teacher)
        teachersList []Teacher
    )
    data := readFile(filename)
    err := json.Unmarshal([]byte(data), &teachersList)
    if err != nil {
        return nil
    }
    for _, tc := range teachersList {
        teachers[tc.Name] = &Teacher{
            Name: tc.Name,
            Cabinets: tc.Cabinets,
        }
    }
    return teachers
}

const (
    OUTDATA = "output/"
)
func (gen *Generator) Template() [][]*Schedule {
    var (
        weekLessons = make([][]*Schedule, 7)
        generator = new(Generator)
    )
    unpackJSON(packJSON(gen), generator)
    file, name := CreateXLSX(OUTDATA + "template.xlsx")
    for i := generator.Day; i < generator.Day+7; i++ {
        weekLessons[i % 7] = generator.Generate(nil)
        if gen.Debug {
            generator.WriteXLSX(
                file,
                name,
                weekLessons[i],
                int(i),
            )
        }
    }
    return weekLessons
}

func (gen *Generator) Generate(template [][]*Schedule) []*Schedule {
    var (
        list   []*Schedule
        templt []*Schedule
        groups = getGroups(gen.Groups)
    )
    if template == nil {
        templt = nil
    } else {
        templt = template[gen.Day]
    }
    for _, group := range groups {
        var (
            schedule = gen.NewSchedule(group.Name)
            subjects = getSubjects(group.Subjects)
            countLessons = new(Subgroup)
        )
        if gen.Day == SUNDAY {
            list = append(list, schedule)
            for _, subject := range subjects {
                saved := gen.Groups[group.Name].Subjects[subject.Name].SaveWeek
                gen.Groups[group.Name].Subjects[subject.Name].WeekLessons.A = saved
                gen.Groups[group.Name].Subjects[subject.Name].WeekLessons.B = saved
            }
            continue
        }
        for _, subject := range subjects {
            switch {
            case gen.haveTheoreticalLessons(subject):
                if gen.Debug {
                    fmt.Println(group.Name, subject.Name, ": not splited THEORETICAL;")
                }
                gen.tryGenerate(ALL, THEORETICAL, group, subject, schedule, countLessons, templt)
            // Практические пары начинаются только после завершения всех теоретических.
            default:
                // Если подгруппа неделимая, тогда провести практику в виде полной пары.
                // Иначе разделить практику на две подгруппы.
                if !gen.withSubgroups(group.Name) {
                    if gen.Debug {
                        fmt.Println(group.Name, subject.Name, ": not splited PRACTICAL;")
                    }
                    gen.tryGenerate(ALL, PRACTICAL, group, subject, schedule, countLessons, templt)
                } else {
                    switch RandSubgroup() {
                    case A:
                        if gen.Debug {
                            fmt.Println(group.Name, subject.Name, ": splited (A -> B);")
                        }
                        gen.tryGenerate(A, PRACTICAL, group, subject, schedule, countLessons, templt)
                        gen.tryGenerate(B, PRACTICAL, group, subject, schedule, countLessons, templt)
                    case B:
                        if gen.Debug {
                            fmt.Println(group.Name, subject.Name, ": splited (B -> A);")
                        }
                        gen.tryGenerate(B, PRACTICAL, group, subject, schedule, countLessons, templt)
                        gen.tryGenerate(A, PRACTICAL, group, subject, schedule, countLessons, templt)
                    }
                }
            }
        }
        list = append(list, schedule)
    }
    gen.Reserved.Teachers = make(map[string][]bool)
    gen.Reserved.Cabinets = make(map[string][]bool)
    gen.Day = (gen.Day + 1) % 7
    return sortSchedule(list)
}

func CreateXLSX(filename string) (*xlsx.File, string) {
    file := xlsx.NewFile()
    _, err := file.AddSheet("Init")
    if err != nil {
        return nil, ""
    }
    err = file.Save(filename)
    if err != nil {
        return nil, ""
    }
    return file, filename
}

func (gen *Generator) WriteXLSX(file *xlsx.File, filename string, schedule []*Schedule, iter int) error {
    const (
        colWidth = 30
        rowHeight = 30
    )

    var (
        MAXCOL = uint(3)
    )

    rowsNext := uint(len(schedule)) / MAXCOL
    if rowsNext == 0 || uint(len(schedule)) % MAXCOL != 0 {
        rowsNext += 1
    }

    var (
        
        colNum = uint(NUM_TABLES + 2)
        
        row = make([]*xlsx.Row, colNum * rowsNext) //  * (rowsNext + 1)
        cell *xlsx.Cell
        dayN = gen.Day
        day = ""
    )

    if dayN == SUNDAY {
        dayN = SATURDAY
    } else {
        dayN -= 1
    }

    switch dayN {
    case SUNDAY: day = "Sunday"
    case MONDAY: day = "Monday"
    case TUESDAY: day = "Tuesday"
    case WEDNESDAY: day = "Wednesday"
    case THURSDAY: day = "Thursday"
    case FRIDAY: day = "Friday"
    case SATURDAY: day = "Saturday"
    }

    sheet, err := file.AddSheet(day + "-" + strconv.Itoa(iter))
    if err != nil {
        return err
    }

    sheet.SetColWidth(2, int(MAXCOL)*3+1, COL_W)

    for r := uint(0); r < rowsNext; r++ {
        for i := uint(0); i < colNum; i++ {
            row[(r*colNum)+i] = sheet.AddRow() // (r*rowsNext)+
            row[(r*colNum)+i].SetHeight(ROW_H)
            cell = row[(r*colNum)+i].AddCell()
            if i == 0 {
                cell.Value = "Пара"
                continue
            }
            cell.Value = strconv.Itoa(int(i-1))
        }
    }

    index := uint(0)
    exit: for r := uint(0); r < rowsNext; r++ {
        for i := uint(0); i < MAXCOL; i++ {
            if uint(len(schedule)) <= index {
                break exit
            }

            savedCell := row[(r*colNum)+0].AddCell()
            savedCell.Value = "Группа " + schedule[index].Group

            cell = row[(r*colNum)+0].AddCell()
            cell = row[(r*colNum)+0].AddCell()

            savedCell.Merge(2, 0)

            cell = row[(r*colNum)+1].AddCell()
            cell.Value = "Предмет"

            cell = row[(r*colNum)+1].AddCell()
            cell.Value = "Преподаватель"

            cell = row[(r*colNum)+1].AddCell()
            cell.Value = "Кабинет"

            for j, trow := range schedule[index].Table {
                cell = row[(r*colNum)+uint(j)+2].AddCell()
                if trow.Subject[A] == trow.Subject[B] {
                    cell.Value = trow.Subject[A]
                } else {
                    if trow.Subject[A] != "" {
                        cell.Value = trow.Subject[A] + " (A)"
                    }
                    if trow.Subject[B] != "" {
                        cell.Value += "\n" + trow.Subject[B] + " (B)"
                    }
                }

                cell = row[(r*colNum)+uint(j)+2].AddCell()
                if trow.Teacher[A] == trow.Teacher[B] {
                    cell.Value = trow.Teacher[A]
                } else {
                    if trow.Teacher[A] != "" {
                        cell.Value = trow.Teacher[A]
                    }
                    if trow.Teacher[B] != "" {
                        cell.Value += "\n" + trow.Teacher[B]
                    }
                }

                sheet.SetColWidth(colWidthForCabinets(int(j)))
                cell = row[(r*colNum)+uint(j)+2].AddCell()
                if trow.Cabinet[A] == trow.Cabinet[B] {
                    cell.Value = trow.Cabinet[A]
                } else {
                    if trow.Cabinet[A] != "" {
                        cell.Value = trow.Cabinet[A]
                    }
                    if trow.Cabinet[B] != "" {
                        cell.Value += "\n" + trow.Cabinet[B]
                    }
                }
            }

            index++
        }
    }

    err = file.Save(filename)
    if err != nil {
        return err
    }

    return nil
}

func RandSubgroup() SubgroupType {
    return SubgroupType(random(0, 1))
}

func Load(filename string) *Generator {
    var generator = new(Generator)
    jsonData := readFile(filename)
    err := json.Unmarshal([]byte(jsonData), generator)
    if err != nil {
        return nil
    }
    return generator
}

func (gen *Generator) Dump(filename string) error {
    return writeJSON(filename, gen)
}
