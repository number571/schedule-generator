package schedule

import (
	"fmt"
	"time"
	"errors"
	"strconv"
	"math/rand"
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

func (gen *Generator) NotHaveLessons(subgroup SubgroupType, subject *Subject) bool {
	switch subgroup {
	case A:
		if subject.Subgroup.A == 0 {
			return true
		}
	case B:
		if subject.Subgroup.B == 0 {
			return true
		}
	case ALL:
		if subject.Subgroup.A == 0 && subject.Subgroup.B == 0 {
			return true
		}
	}
	return false
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
			groups[gr.Name].Subjects[sb.Name] = &Subject{
				Name: sb.Name,
				Teacher: sb.Teacher,
				IsSplited: sb.IsSplited,
				All: sb.Lessons.All,
				Subgroup: Subgroup{
					A: sb.Lessons.WeekLessons,
					B: sb.Lessons.WeekLessons,
				},
			}
		}
	}
	return groups
}

func ReadTeachers(filename string) map[string]*Teacher {
	var teachers = make(map[string]*Teacher)
	data := readFile(filename)
	err := json.Unmarshal([]byte(data), &teachers)
	if err != nil {
		return nil
	}
	return teachers
}

func (gen *Generator) Generate() []*Schedule {
	var list []*Schedule
	groups := getGroups(gen.Groups)
	for _, group := range groups {
		var schedule = gen.NewSchedule(group.Name)
		if gen.Day == SUNDAY {
			list = append(list, schedule)
			continue
		}
		subjects := getSubjects(group.Subjects)
		for _, subject := range subjects {
			if !subject.IsSplited {
				if DEBUG {
					fmt.Println(group.Name, subject.Name, ": not splited;")
				}
				gen.tryGenerate(ALL, group, subject, schedule)
			} else {
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

func WriteXLSX(file *xlsx.File, filename string, schedule []*Schedule, numtable uint8, iter int) error {
    const (
        colWidth = 30
        rowHeight = 30
    )

    var (
        colNum = numtable + 1
        row = make([]*xlsx.Row, colNum)
        cell *xlsx.Cell
    )

    sheet, err := file.AddSheet("Schedule-" + strconv.Itoa(iter))
    if err != nil {
        return err
    }

    sheet.SetColWidth(2, len(schedule)*3+1, COL_W)
    for i := uint8(0); i < colNum; i++ {
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

        for j, trow := range sch.Table {

            cell = row[j+1].AddCell()
            if trow.Teacher[0] == trow.Teacher[1] {
                cell.Value = trow.Teacher[0]
            } else {
                if trow.Teacher[0] != "" {
                    cell.Value = trow.Teacher[0]
                }
                if trow.Teacher[1] != "" {
                    cell.Value += "\n" + trow.Teacher[1]
                }
            }
            
            cell = row[j+1].AddCell()
            if trow.Subject[0] == trow.Subject[1] {
                cell.Value = trow.Teacher[0]
            } else {
                if trow.Subject[0] != "" {
                    cell.Value = trow.Subject[0] + " (A)"
                }
                if trow.Subject[1] != "" {
                    cell.Value += "\n" + trow.Subject[1] + " (B)"
                }
            }

            sheet.SetColWidth(colWidthForCabinets(i))
            cell = row[j+1].AddCell()
            if trow.Cabinet[0] == trow.Cabinet[1] {
                cell.Value = trow.Cabinet[0]
            } else {
                if trow.Cabinet[0] != "" {
                    cell.Value = trow.Cabinet[0]
                }
                if trow.Cabinet[1] != "" {
                    cell.Value += "\n" + trow.Cabinet[1]
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

func Shuffle(slice interface{}) interface{}{
    switch slice.(type) {
    case []*Group:
        result := slice.([]*Group)
        rand.Seed(int64(time.Now().Nanosecond()))
        for i := len(result)-1; i > 0; i-- {
            j := rand.Intn(i+1)
            result[i], result[j] = result[j], result[i]
        }
        return result
    case []*Subject:
        result := slice.([]*Subject)
        rand.Seed(int64(time.Now().Nanosecond()))
        for i := len(result)-1; i > 0; i-- {
            j := rand.Intn(i+1)
            result[i], result[j] = result[j], result[i]
        }
        return result
    }
    return nil
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

func (gen *Generator) UnblockTeacher(teacher string) error {
	if !gen.InBlocked(teacher) {
		return errors.New("teacher undefined")
	}
	delete(gen.Blocked, teacher)
	return nil
}

func WriteJSON(filename string, data interface{}) error {
	return writeFile(filename, string(packJSON(data)))
}

func packJSON(data interface{}) []byte {
	jsonData, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return nil
	}
	return jsonData
}

func RandSubgroup() SubgroupType {
	return SubgroupType(random(0, 1))
}

// Returns [min:max] value.
func random(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max - min + 1) + min
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

func (gen *Generator) InBlocked(teacher string) bool {
	if _, ok := gen.Blocked[teacher]; ok {
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

func (gen *Generator) Dump(filename string) error {
	return WriteJSON(filename, gen)
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


func (gen *Generator) BlockTeacher(teacher string) error {
	if !gen.InTeachers(teacher) {
		return errors.New("teacher undefined")
	}
	gen.Blocked[teacher] = true
	return nil
}
