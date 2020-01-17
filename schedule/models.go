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
	Semester uint8
	NumTables uint8
	Groups map[string]*Group
	Teachers map[string]Teacher
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
	Teacher [2]string
	Subject [2]string
	Cabinet [2]string
}

type Teacher struct {
	Cabinets []string
	Groups map[string]string
}

// QWEasd123

// type Subject struct {
// 	Teachers map[string]bool
// 	Groups map[string]bool
// }

// type Cabinet struct {
// 	IsComputer bool
// 	Places uint16 // free places in cabinet
// 	Teachers map[string]bool
// 	Subjects map[string]bool
// }

type Group struct {
	Name string
	Quantity uint16 // students count
	Subjects map[string]*Subject
}

type Subject struct {
	Name string
	Teacher string
	IsSplited bool
	Subgroup Subgroup
}

type Subgroup struct {
	A [2]Semester
	B [2]Semester
}

// type Hours struct {
// 	All uint16
// 	Semester []Semester
// }

type Semester struct {
	All uint16
	WeekHours uint8
}
