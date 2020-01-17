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
	PrintJSON(generator.Generate())
	// generator.BlockTeacher("teacher_two")
	PrintJSON(generator.Generate())
	// generator.UnblockTeacher("teacher_two")
	// PrintJSON(generator.Generate())
}

func PrintJSON(data interface{}) {
	jsonData, _ := json.MarshalIndent(data, "", "\t")
	fmt.Println(string(jsonData))
}

func NewSubject(subject string, teacher string, splited bool) *schedule.Subject {
	return &schedule.Subject{
		Name: subject,
		Teacher: teacher,
		IsSplited: splited,
		Subgroup: schedule.Subgroup{
			A: [2]schedule.Semester{
				schedule.Semester{
					All: 28,
					WeekHours: 4, // 2 couples
				},
				schedule.Semester{
					All: 22,
					WeekHours: 2,
				},
			},
			B: [2]schedule.Semester{
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

var groups = map[string]*schedule.Group{
	"101": &schedule.Group{
		Name: "101",
		Quantity: 20,
		Subjects: map[string]*schedule.Subject{
			"subject_one": NewSubject("subject_one", "teacher_one", true),
			"subject_two": NewSubject("subject_two", "teacher_two", false),
		},
	},
	"102": &schedule.Group{
		Name: "102",
		Quantity: 15,
		Subjects: map[string]*schedule.Subject{
			"subject_one": NewSubject("subject_one", "teacher_one", true),
			"subject_two": NewSubject("subject_two", "teacher_two", false),
		},
	},
}
