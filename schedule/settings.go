package schedule

var Teachers = map[string]Teacher{
	"teacher_one": Teacher{
		Subjects: map[string]bool{
			"subject_one": true,
		},
		Cabinets: map[string]bool{
			"201": true,
		},
		Groups: map[string]bool{
			"101": true,
			"102": true,
			"103": true,
		},
	},
	"teacher_two": Teacher{
		Subjects: map[string]bool{
			"subject_two": true,
		},
		Cabinets: map[string]bool{
			"301": true,
		},
		Groups: map[string]bool{
			"101": true,
			"102": true,
			"103": true,
		},
	},
}

var Groups = map[string]Group{
	"101": Group{
		Quantity: 20,
		Subjects: map[string]LocalSubj{
			"subject_one": LocalSubj{
				Teacher: "teacher_one",
				Hours: Hours{
					All: 50,
					Semester: []Semester{
						Semester{
							All: 28,
							WeekHours: 2,
						},
						Semester{
							All: 22,
							WeekHours: 4,
						},
					},
				},
			},
			"subject_two": LocalSubj{
				Teacher: "teacher_two",
				Hours: Hours{
					All: 50,
					Semester: []Semester{
						Semester{
							All: 28,
							WeekHours: 2,
						},
						Semester{
							All: 22,
							WeekHours: 4,
						},
					},
				},
			},
			// 
		},
	},

	"102": Group{
		Quantity: 15,
		Subjects: map[string]LocalSubj{
			"subject_one": LocalSubj{
				Teacher: "teacher_one",
				Hours: Hours{
					All: 50,
					Semester: []Semester{
						Semester{
							All: 28,
							WeekHours: 2,
						},
						Semester{
							All: 22,
							WeekHours: 4,
						},
					},
				},
			},
			"subject_two": LocalSubj{
				Teacher: "teacher_two",
				Hours: Hours{
					All: 50,
					Semester: []Semester{
						Semester{
							All: 28,
							WeekHours: 2,
						},
						Semester{
							All: 22,
							WeekHours: 4,
						},
					},
				},
			},
			// 
		},
	},

	"103": Group{
		Quantity: 18,
		Subjects: map[string]LocalSubj{
			"subject_one": LocalSubj{
				Teacher: "teacher_one",
				Hours: Hours{
					All: 50,
					Semester: []Semester{
						Semester{
							All: 28,
							WeekHours: 2,
						},
						Semester{
							All: 22,
							WeekHours: 4,
						},
					},
				},
			},
			"subject_two": LocalSubj{
				Teacher: "teacher_two",
				Hours: Hours{
					All: 50,
					Semester: []Semester{
						Semester{
							All: 28,
							WeekHours: 2,
						},
						Semester{
							All: 22,
							WeekHours: 4,
						},
					},
				},
			},
			// 
		},
	},

}
