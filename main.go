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
	result := generator.Generate()
	PrintJSON(result)
	generator.WriteXLSX("schedule.xlsx", result)
}

func PrintJSON(data interface{}) {
	jsonData, _ := json.MarshalIndent(data, "", "\t")
	fmt.Println(string(jsonData))
}

func NewSubject(name string, teacher string, splited bool, semesters [2]schedule.Semester) *schedule.Subject {
	return &schedule.Subject{
		Name: name,
		Teacher: teacher,
		IsSplited: splited,
		Subgroup: schedule.Subgroup{
			A: [2]schedule.Semester{
				schedule.Semester{
					All: semesters[0].All,
					WeekLessons: semesters[0].WeekLessons,
				},
				schedule.Semester{
					All: semesters[1].All,
					WeekLessons: semesters[1].WeekLessons,
				},
			},
			B: [2]schedule.Semester{
				schedule.Semester{
					All: semesters[0].All,
					WeekLessons: semesters[0].WeekLessons,
				},
				schedule.Semester{
					All: semesters[1].All,
					WeekLessons: semesters[1].WeekLessons,
				},
			},
		},
	}
}

var semesterExample = [2]schedule.Semester{
	schedule.Semester{
		All: 14,
		WeekLessons: 2, // 2 lessons
	},
	schedule.Semester{
		All: 11,
		WeekLessons: 1, // 1 lesson
	},
}

var teachers = map[string]*schedule.Teacher{
	"teacher_one": &schedule.Teacher{
		Cabinets: []string{
			"201", 
			"204",
		},
		Groups: map[string]string{
			"101": "subject_one",
			"102": "subject_one",
			"103": "subject_one",
			"104": "subject_one",
			"105": "subject_one",
		},
	},
	"teacher_two": &schedule.Teacher{
		Cabinets: []string{
			"301",
		},
		Groups: map[string]string{
			"101": "subject_two",
			"102": "subject_two",
			"103": "subject_two",
			"104": "subject_two",
			"105": "subject_two",
		},
	},
	"teacher_three": &schedule.Teacher{
		Cabinets: []string{
			"304",
		},
		Groups: map[string]string{
			"101": "subject_three",
			"102": "subject_three",
			"103": "subject_three",
			"104": "subject_three",
			"105": "subject_three",
		},
	},
}

var groups = map[string]*schedule.Group{
	"101": &schedule.Group{
		Name: "101",
		Quantity: 20,
		Subjects: map[string]*schedule.Subject{
			"subject_one": NewSubject("subject_one", "teacher_one", true, semesterExample),
			"subject_two": NewSubject("subject_two", "teacher_two", false, semesterExample),
			"subject_three": NewSubject("subject_three", "teacher_three", false, semesterExample),
		},
	},
	"102": &schedule.Group{
		Name: "102",
		Quantity: 15,
		Subjects: map[string]*schedule.Subject{
			"subject_one": NewSubject("subject_one", "teacher_one", false, semesterExample),
			"subject_two": NewSubject("subject_two", "teacher_two", true, semesterExample),
			"subject_three": NewSubject("subject_three", "teacher_three", false, semesterExample),
		},
	},
	"103": &schedule.Group{
		Name: "103",
		Quantity: 15,
		Subjects: map[string]*schedule.Subject{
			"subject_one": NewSubject("subject_one", "teacher_one", true, semesterExample),
			"subject_two": NewSubject("subject_two", "teacher_two", false, semesterExample),
			"subject_three": NewSubject("subject_three", "teacher_three", false, semesterExample),
		},
	},
	"104": &schedule.Group{
		Name: "104",
		Quantity: 15,
		Subjects: map[string]*schedule.Subject{
			"subject_one": NewSubject("subject_one", "teacher_one", true, semesterExample),
			"subject_two": NewSubject("subject_two", "teacher_two", false, semesterExample),
			"subject_three": NewSubject("subject_three", "teacher_three", false, semesterExample),
		},
	},
	"105": &schedule.Group{
		Name: "105",
		Quantity: 15,
		Subjects: map[string]*schedule.Subject{
			"subject_one": NewSubject("subject_one", "teacher_one", false, semesterExample),
			"subject_two": NewSubject("subject_two", "teacher_two", false, semesterExample),
			"subject_three": NewSubject("subject_three", "teacher_three", true, semesterExample),
		},
	},
}
