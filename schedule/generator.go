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
        teachersList []TeacherJSON
    )
    data := readFile(filename)
    err := json.Unmarshal([]byte(data), &teachersList)
    if err != nil {
        return nil
    }
    for _, tc := range teachersList {
        teachers[tc.Name] = &Teacher{
            Cabinets: tc.Cabinets,
        }
    }
    return teachers
}

func printJSON(data interface{}) {
    jsonData, _ := json.MarshalIndent(data, "", "\t")
    fmt.Println(string(jsonData))
}

func (gen *Generator) Generate() []*Schedule {
    var list []*Schedule
    groups := getGroups(gen.Groups)
    for _, group := range groups {
        var schedule = gen.NewSchedule(group.Name)
        subjects := getSubjects(group.Subjects)
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
            if DEBUG {
                fmt.Println(group.Name, subject.Name, ": not splited;")
            }
            gen.tryGenerate(ALL, group, subject, schedule)

            // Практические пары начинаются только после завершения всех теоретических.
            if !gen.haveTheoreticalLessons(subject) {
                switch RandSubgroup() {
                case A:
                    if DEBUG {
                        fmt.Println(group.Name, subject.Name, ": splited (A -> B);")
                    }
                    gen.tryGenerate(A, group, subject, schedule)
                    gen.tryGenerate(B, group, subject, schedule)
                case B:
                    if DEBUG {
                        fmt.Println(group.Name, subject.Name, ": splited (B -> A);")
                    }
                    gen.tryGenerate(B, group, subject, schedule)
                    gen.tryGenerate(A, group, subject, schedule)
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
        colNum = gen.NumTables + 2
        row = make([]*xlsx.Row, colNum)
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

    sheet.SetColWidth(2, len(schedule)*3+1, COL_W)
    for i := uint(0); i < colNum; i++ {
        row[i] = sheet.AddRow()
        row[i].SetHeight(ROW_H)
        cell = row[i].AddCell()
        if i == 0 {
            cell.Value = "Пара"
            continue
        }
        cell.Value = strconv.Itoa(int(i))
    }

    for i, sch := range schedule {
        savedCell := row[0].AddCell()
        savedCell.Value = "Группа " + sch.Group

        cell = row[0].AddCell()
        cell = row[0].AddCell()

        savedCell.Merge(2, 0)

        cell = row[1].AddCell()
        cell.Value = "Предмет"

        cell = row[1].AddCell()
        cell.Value = "Преподаватель"

        cell = row[1].AddCell()
        cell.Value = "Кабинет"

        for j, trow := range sch.Table {

            cell = row[j+2].AddCell()
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

            cell = row[j+2].AddCell()
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

            sheet.SetColWidth(colWidthForCabinets(i))
            cell = row[j+2].AddCell()
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
