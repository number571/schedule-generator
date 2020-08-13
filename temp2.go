package main

import (
	"os"
	"fmt"
	"bytes"
	"strings"
	"net/http"
	"io/ioutil"
)

var (
	teachers = readFile("input/teachers.json")
	groups = readFile("input/groups.json")
	generator = readFile("input/generator.json")
)

func main() {
	var fmtstr = ""

	if len(os.Args) != 2 {
		panic("args != 2")
		return
	}

	switch {
	case strings.Contains(os.Args[1], "create"):
		fmtstr = fmt.Sprintf("{\"day\": %d, \"teachers\": %s, \"groups\": %s}", 1, teachers, groups)
	case strings.Contains(os.Args[1], "generate"):
		fmtstr = fmt.Sprintf("{\"generator\": %s}", generator)
	}

	resp, err := http.Post(
		os.Args[1],
		"application/json",
		bytes.NewBuffer([]byte(fmtstr)),
	)
	if err != nil {
		panic(err)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))
}

func readFile(filename string) string {
    data, err := ioutil.ReadFile(filename)
    if err != nil {
    	return ""
    }
    return string(data)
}

