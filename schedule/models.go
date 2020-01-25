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

type Generator struct {
	Day DayType
	NumTables uint8
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
	Subject [2]string
	Teacher [2]string
	Cabinet [2]string
}

type Teacher struct {
	Cabinets []string
	// Groups map[string]string
}

type Group struct {
	Name string
	Quantity uint16 // students count
	Subjects map[string]*Subject
}

type Subject struct {
	Name string
	Teacher string
	Teacher2 string
	IsSplited bool
	All uint16
	Subgroup Subgroup
}

type Subgroup struct {
	A uint8
	B uint8
}

type TeacherJSON struct {
	Name string `json:"name"`
	Cabinets []string `json:"cabinets"`
}

type GroupJSON struct {
	Name string `json:"name"`
	Quantity uint16 `json:"quantity"`
	Subjects []SubjectJSON `json:"subjects"`
}

type SubjectJSON struct {
	Name string `json:"name"`
	Teacher string `json:"teacher"`
	IsSplited bool `json:"is_splited"`
	Lessons LessonsJSON `json:"lessons"`
}

type LessonsJSON struct {
	All uint16 `json:"all"`
	Week uint8 `json:"week"`
}
