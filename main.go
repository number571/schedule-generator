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
		Debug: true,
		Day: schedule.MONDAY,
		Groups: schedule.ReadGroups(INDATA + "groups.json"),
		Teachers: schedule.ReadTeachers(INDATA + "teachers.json"),
	})
	template := generator.Template()
	os.Mkdir(OUTDATA, 0777)
	file, name := schedule.CreateXLSX(OUTDATA + "schedule.xlsx")
	for iter := 1; iter <= 7; iter++ {
		result := generator.Generate(template)
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

/*
// Groups:
[
	{
		"name": "101",
		"quantity": 20,
		"subjects": [
			{
				"name": "subject_one",
				"teacher": "teacher_one",
				"is_computer": true,
				"lessons": {
					"theory": 2,
					"practice": 2,
					"week": 6
				}
			}
		]
	}
]

// Teachers:
[
	{
		"name": "teacher_one",
		"cabinets": [
			{
				"name": "201",
				"is_computer": true
			}
		]
	}
]
*/
