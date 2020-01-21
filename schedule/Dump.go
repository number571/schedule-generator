package schedule

import (
	"os"
	"encoding/json"
)

func (gen *Generator) Dump(filename string) error {
	jsonData, _ := json.MarshalIndent(gen, "", "\t")
	return writeFile(filename, string(jsonData))
}

func writeFile(filename string, data string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	file.WriteString(data)
	file.Close()
	return nil
}
