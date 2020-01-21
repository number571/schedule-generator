package schedule

func (gen *Generator) NewSchedule(group string) *Schedule {
	return &Schedule{
		Day: gen.Day,
		Group: group,
		Table: make([]Row, gen.NumTables),
	}
}
