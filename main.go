package main

import (
	"fmt"
	"encoding/json"
	"./schedule"
)

func main() {
	var generator = schedule.NewGenerator(&schedule.Generator{
		Day: schedule.MONDAY,
		Semester: 1,
		NumTables: 11,
		Groups: groups,
		Teachers: teachers,
	})
	PrintJSON(generator.Generate(condition))
	// generator.BlockTeacher("teacher_two")
	// PrintJSON(generator.Generate(condition))
	// generator.UnblockTeacher("teacher_two")
	// PrintJSON(generator.Generate(condition))
}

// true = continue;
// false = break group schedule;
func condition(gen *schedule.Generator, grp *schedule.Group, sbj *schedule.Subject, sch *schedule.Schedule, cpl *uint8) bool {
	// Full day = max 6 couples.
	if (gen.Day != schedule.WEDNESDAY && gen.Day != schedule.SATURDAY) && *cpl == 6 {
		return false
	}
	// Without middle spaces.
	if *cpl > 1 && gen.IsReserved(sch, *cpl-2) && !gen.IsReserved(sch, *cpl-1) {
		return false
	}
	return true
}

func PrintJSON(data interface{}) {
	jsonData, _ := json.MarshalIndent(data, "", "\t")
	fmt.Println(string(jsonData))
}

func NewSubject(subject string, teacher string) *schedule.Subject {
	return &schedule.Subject{
		Name: subject,
		Teacher: teacher,
		Hours: schedule.Hours{
			All: 50,
			Semester: []schedule.Semester{
				schedule.Semester{
					All: 28,
					WeekHours: 4, // 2 couples
				},
				schedule.Semester{
					All: 22,
					WeekHours: 2,
				},
			},
		},
	}
}

var teachers = map[string]schedule.Teacher{
	"teacher_one": schedule.Teacher{
		Cabinets: []string{
			"201", 
			"204",
		},
		Groups: map[string]string{
			"101": "subject_one",
			"102": "subject_one",
			"103": "subject_one",
		},
	},
	"teacher_two": schedule.Teacher{
		Cabinets: []string{
			"301",
		},
		Groups: map[string]string{
			"101": "subject_two",
			"102": "subject_two",
			"103": "subject_two",
		},
	},
}

var groups = map[string]schedule.Group{
	"101": schedule.Group{
		Quantity: 20,
		Subjects: map[string]*schedule.Subject{
			"subject_one": NewSubject("subject_one", "teacher_one"),
			"subject_two": NewSubject("subject_two", "teacher_two"),
		},
	},
	"102": schedule.Group{
		Quantity: 15,
		Subjects: map[string]*schedule.Subject{
			"subject_one": NewSubject("subject_one", "teacher_one"),
			"subject_two": NewSubject("subject_two", "teacher_two"),
		},
	},
	"103": schedule.Group{
		Quantity: 18,
		Subjects: map[string]*schedule.Subject{
			"subject_one": NewSubject("subject_one", "teacher_one"),
			"subject_two": NewSubject("subject_two", "teacher_two"),
		},
	},
}
