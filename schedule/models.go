package schedule

type DayType uint8
const (
	SUNDAY 		DayType = 0
	MONDAY 		DayType = 1
	TUESDAY 	DayType = 2
	WEDNESDAY 	DayType = 3
	THURSDAY 	DayType = 4
	FRIDAY 		DayType = 5
	SATURDAY 	DayType = 6
)

type SubgroupType uint8
const (
	A 	SubgroupType = 0
	B 	SubgroupType = 1
	ALL SubgroupType = 2
)

type SubjectType uint8
const (
	THEORETICAL	SubjectType = 0
	PRACTICAL 	SubjectType = 1
)

type Generator struct {
	Day DayType
	NumTables uint
	Groups map[string]*Group
	Teachers map[string]*Teacher
	Blocked map[string]bool
	Reserved Reserved
}

type Reserved struct {
	Teachers map[string][]bool
	Cabinets map[string][]bool
}

type Schedule struct {
	Day DayType
	Group string
	Table []Row
}

type Row struct {
	Subject [ALL]string
	Teacher [ALL]string
	Cabinet [ALL]string
}

type Teacher struct {
	Cabinets []string
}

type Group struct {
	Name string
	Quantity uint // students count
	Subjects map[string]*Subject
}

type Subject struct {
	Name string
	Teacher string
	Teacher2 string
	SaveWeek uint
	Theory uint
	Practice Subgroup
	WeekLessons Subgroup
}

type Subgroup struct {
	A uint
	B uint
}

type TeacherJSON struct {
	Name string `json:"name"`
	Cabinets []string `json:"cabinets"`
}

type GroupJSON struct {
	Name string `json:"name"`
	Quantity uint `json:"quantity"`
	Subjects []SubjectJSON `json:"subjects"`
}

type SubjectJSON struct {
	Name string `json:"name"`
	Teacher string `json:"teacher"`
	Lessons LessonsJSON `json:"lessons"`
}

type LessonsJSON struct {
	Theory uint `json:"theory"`
	Practice uint `json:"practice"`
	Week uint `json:"week"`
}
