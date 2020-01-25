package main

/*
	github.com/tealeg/xlsx
*/

import (
	"os"
	"./schedule"
)

const (
	INDATA = "input/"
	OUTDATA = "output/"
	XLSX = "schedule.xlsx"
)

func main() {
	var generator = schedule.NewGenerator(&schedule.Generator{
		Day: schedule.MONDAY,
		NumTables: 11,
		Groups: schedule.ReadGroups(INDATA + "groups.json"),
		Teachers: schedule.ReadTeachers(INDATA + "teachers.json"),
	})
	os.Mkdir(OUTDATA, 0777)
	file, name := schedule.CreateXLSX(OUTDATA + XLSX)
	for iter := 1; iter <= 3; iter++ {
		result := generator.Generate()
		schedule.WriteXLSX(
			file,
			name,
			result,
			generator.NumTables,
			iter,
		)
	}
}
