package schedule

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
