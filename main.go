package main

/*
	github.com/tealeg/xlsx
*/

import (
	"os"
	"fmt"
	"encoding/json"
	"./schedule"
)

const (
	INDATA = "input/"
	OUTDATA = "output/"
)

func main() {
	var generator = schedule.NewGenerator(&schedule.Generator{
		Day: schedule.FRIDAY,
		Groups: schedule.ReadGroups(INDATA + "groups.json"),
		Teachers: schedule.ReadTeachers(INDATA + "teachers.json"),
	})
	os.Mkdir(OUTDATA, 0777)
	file, name := schedule.CreateXLSX(OUTDATA + "schedule.xlsx")
	for iter := 1; iter <= 15; iter++ {
		result := generator.Generate()
		generator.WriteXLSX(
			file,
			name,
			result,
			iter,
		)
	}
	// printJSON(generator)
}

func printJSON(data interface{}) {
	jsonData, _ := json.MarshalIndent(data, "", "\t")
	fmt.Println(string(jsonData))
}
