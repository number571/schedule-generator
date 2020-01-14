package main

import (
	"fmt"
	"encoding/json"
	"./schedule"
)

func main() {
	generator := schedule.NewGenerator(&schedule.GenData{
		Day: schedule.MONDAY,
		Semester: 0,
		NumTables: 11,
		Groups: schedule.Groups,
	})
	PrintJSON(generator.Generate())
}

func PrintJSON(data interface{}) {
	jsonData, _ := json.MarshalIndent(data, "", "\t")
	fmt.Println(string(jsonData))
}
