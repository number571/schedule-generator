package schedule

import (
	"os"
	"encoding/json"
)

func Load(filename string) *Generator {
	var generator Generator
	jsonData := readFile(filename)
	err := json.Unmarshal([]byte(jsonData), &generator)
	if err != nil {
		return nil
	}
	return &generator
}

func readFile(filename string) string {
	file, err := os.Open(filename)
	if err != nil {
		return ""
	}
	defer file.Close()

	var (
		buffer []byte = make([]byte, 512)
		data string
	)

	for {
		length, err := file.Read(buffer)
		if length == 0 || err != nil { break }
		data += string(buffer[:length])
	}

	return data
}
