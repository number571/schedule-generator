package schedule

import (
	"fmt"
	"sort"
)

func (gen *Generator) Generate() []*Schedule {
	var list []*Schedule
	groups := getGroups(gen.Groups)
	for _, group := range groups {
		var schedule = gen.NewSchedule(group.Name)
		if gen.Day == SUNDAY {
			list = append(list, schedule)
			continue
		}
		subjects := getSubjects(group.Subjects)
		for _, subject := range subjects {
			if !subject.IsSplited {
				if DEBUG {
					fmt.Println(group.Name, subject.Name, ": not splited;")
				}
				gen.tryGenerate(ALL, group, subject, schedule)
			} else {
				switch RandSubgroup() {
				case A:
					if DEBUG {
						fmt.Println(group.Name, subject.Name, ": splited (A -> B);")
					}
					gen.tryGenerate(A, group, subject, schedule)
					gen.tryGenerate(B, group, subject, schedule)
				case B:
					if DEBUG {
						fmt.Println(group.Name, subject.Name, ": splited (B -> A);")
					}
					gen.tryGenerate(B, group, subject, schedule)
					gen.tryGenerate(A, group, subject, schedule)
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

func getGroups(groups map[string]*Group) []*Group {
	var list []*Group
	for _, group := range groups {
		list = append(list, group)
	}
	return Shuffle(list).([]*Group)
}

func getSubjects(subjects map[string]*Subject) []*Subject {
	var list []*Subject
	for _, subject := range subjects {
		list = append(list, subject)
	}
	return Shuffle(list).([]*Subject)
}

func sortSchedule(schedule []*Schedule) []*Schedule {
	sort.SliceStable(schedule, func(i, j int) bool {
		return schedule[i].Group < schedule[j].Group
	})
	return schedule
}
