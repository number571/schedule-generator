package schedule

import (
	"fmt"
	"sort"
)

func (gen *Generator) Generate() []*Schedule {
	var list []*Schedule
	for _, group := range gen.Groups {
		var schedule = gen.NewSchedule(group.Name)
		if gen.Day == SUNDAY {
			list = append(list, schedule)
			continue
		}
		nextsub: for _, subject := range group.Subjects {
			if !subject.IsSplited {
				fmt.Println(group.Name, ": not splited;")
				if gen.tryGenerate(ALL, group, subject, schedule) {
					break nextsub
				}
			} else {
				fmt.Println(group.Name, ": splited;")
				if gen.tryGenerate(A, group, subject, schedule) {
					break nextsub
				}
				if gen.tryGenerate(B, group, subject, schedule) {
					break nextsub
				}
			}
		}
		list = append(list, schedule)
	}
	gen.Reserved.Teachers = make(map[string][]bool)
	gen.Reserved.Cabinets = make(map[string][]bool)
	gen.Day = (gen.Day + 1) % 7
	return sortSchedule(list)
}

func sortSchedule(schedule []*Schedule) []*Schedule {
	sort.SliceStable(schedule, func(i, j int) bool {
		return schedule[i].Group < schedule[j].Group
	})
	return schedule
}
