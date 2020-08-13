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
	reserved reserved
	Day DayType `json:"day"`
	Groups map[string]*Group `json:"groups"`
	Teachers map[string]*Teacher `json:"teachers"`
	Blocked Blocked `json:"blocked"`
}

type Blocked struct {
	Groups []string `json:"groups"`
	Teachers []string `json:"teachers"`
}

type reserved struct {
	teachers map[string][]bool
	cabinets map[string][]bool
}

type Schedule struct {
	Day DayType `json:"day"`
	Group string `json:"group"`
	Table []Row `json:"table"`
}

type Row struct {
	Subject [ALL]string `json:"subject"`
	Teacher [ALL]string `json:"teacher"`
	Cabinet [ALL]string `json:"cabinet"`
}

type Teacher struct {
	Name string `json:"name"`
	Cabinets []Cabinet `json:"cabinets"`
}

type Cabinet struct {
	Name string `json:"name"`
	IsComputer bool `json:"is_computer"`
}

type Group struct {
	Name string `json:"name"`
	Quantity uint `json:"quantity"` // students count
	Subjects map[string]*Subject `json:"subjects"`
}

type Subject struct {
	Name string `json:"name"`
	Teacher string `json:"teacher"`
	Teacher2 string `json:"teacher2"`
	IsComputer bool `json:"is_computer"`
	SaveWeek uint `json:"save_week"`
	Theory uint `json:"theory"`
	Practice Subgroup `json:"practice"`
	WeekLessons Subgroup `json:"week_lessons"`
}

type Subgroup struct {
	A uint `json:"a"`
	B uint `json:"b"`
}

type GroupJSON struct {
	Name string `json:"name"`
	Quantity uint `json:"quantity"`
	Subjects []SubjectJSON `json:"subjects"`
}

type SubjectJSON struct {
	Name string `json:"name"`
	Teacher string `json:"teacher"`
	IsComputer bool `json:"is_computer"`
	Lessons LessonsJSON `json:"lessons"`
}

type LessonsJSON struct {
	Theory uint `json:"theory"`
	Practice uint `json:"practice"`
	Week uint `json:"week"`
}
