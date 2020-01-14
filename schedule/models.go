package schedule

type Day uint8
const (
	SUNDAY 		Day = 0
	MONDAY 		Day = 1
	TUESDAY 	Day = 2
	WEDNESDAY 	Day = 3
	THURSDAY 	Day = 4
	FRIDAY 		Day = 5
	SATURDAY 	Day = 6
)

type GenData struct {
	Day Day
	Semester uint8
	NumTables uint8
	Groups map[string]Group
}

type Generator struct {
	Day Day
	Semester uint8
	NumTables uint8
	Groups map[string]Group
	Reserved Reserved
}

type Reserved struct {
	Teachers map[string][]bool
	Cabinets map[string][]bool
}

type Schedule struct {
	Day Day
	Group string
	Table []Row
}

type Row struct {
	Teacher string
	Subject string
	Cabinet string
}

type Teacher struct {
	Subjects map[string]bool
	Cabinets map[string]bool
	Groups map[string]bool
}

type Subject struct {
	Teachers map[string]bool
	Groups map[string]bool
}

type Cabinet struct {
	IsComputer bool
	Places uint16 // free places in cabinet
	Teachers map[string]bool
	Subjects map[string]bool
}

type Group struct {
	Quantity uint16 // students count
	Subjects map[string]LocalSubj
}

type LocalSubj struct {
	Teacher string
	Hours Hours
}

type Hours struct {
	All uint16
	Semester []Semester
}

type Semester struct {
	All uint16
	WeekHours uint8
}
