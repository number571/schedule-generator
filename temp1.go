package main

/*
	github.com/tealeg/xlsx
*/

import (
	"os"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"./schedule"
)

const (
	INDATA = "input/"
	OUTDATA = "output/"
)

func main() {
	var (
		groups []schedule.GroupJSON
		teachers []schedule.Teacher
	)
	json.Unmarshal(readFile(INDATA + "groups.json"), &groups)
	json.Unmarshal(readFile(INDATA + "teachers.json"), &teachers)
	var generator = schedule.NewGenerator(&schedule.Generator{
		Day: schedule.MONDAY,
		Groups: schedule.ReadGroups(groups),
		Teachers: schedule.ReadTeachers(teachers),
	})
	template := generator.Template()
	os.Mkdir(OUTDATA, 0777)
	file, name := schedule.CreateXLSX(OUTDATA + "schedule.xlsx")
	// generator.Generate(template)
	for iter := 1; iter <= 6; iter++ {
		// generator.Generate(template)
		result := generator.Generate(template)
		// printJSON(result)
		generator.WriteXLSX(
			file,
			name,
			result,
			iter,
		)
	}
	printJSON(generator)
}

func printJSON(data interface{}) {
	jsonData, _ := json.MarshalIndent(data, "", "\t")
	fmt.Println(string(jsonData))
}

func readFile(filename string) []byte {
    data, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil
    }
    return data
}
