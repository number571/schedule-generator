package main

import (
	"net/http"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	sc "./schedule"
)

var (
	Template [][]*sc.Schedule
)

func main() {
	http.HandleFunc("/", indexPage)
	http.HandleFunc("/create", createPage)
	http.HandleFunc("/generate", generatePage)
	http.ListenAndServe(":7545", nil)
}

func indexPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct{
		Return int `json:"return"`
	}{
		Return: 0,
	})
}

// Method: POST
func createPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var data struct {
		Generator *sc.Generator `json:"generator"`
		Return int `json:"return"`
	}

	var read struct {
		Day int `json:"day"`
		Groups []sc.GroupJSON `json:"groups"`
		Teachers []sc.Teacher `json:"teachers"`
	}

	if r.Method != "POST" {
		data.Return = 1
		json.NewEncoder(w).Encode(data)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&read)
	if err != nil {
		data.Return = 2
		json.NewEncoder(w).Encode(data)
		return
	}

	groups := sc.ReadGroups(read.Groups)
	if groups == nil {
		data.Return = 3
		json.NewEncoder(w).Encode(data)
		return
	}

	teachers := sc.ReadTeachers(read.Teachers)
	if teachers == nil {
		data.Return = 4
		json.NewEncoder(w).Encode(data)
		return
	}

	generator := sc.NewGenerator(&sc.Generator{
		Day: sc.DayType(read.Day),
		Groups: groups,
		Teachers: teachers,
	})

	Template = generator.Template()
	data.Generator = generator

	enc := json.NewEncoder(w)
	enc.SetIndent("", "\t")
	enc.Encode(data)
}

// Method: POST
func generatePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var data struct {
		Schedule []*sc.Schedule `json:"schedule"`
		Hashsum string `json:"hashsum"`
		Return int `json:"return"`
	}
	
	var read struct {
		Generator *sc.Generator `json:"generator"`
	}

	if r.Method != "POST" {
		data.Return = 1
		json.NewEncoder(w).Encode(data)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&read)
	if err != nil {
		data.Return = 2
		json.NewEncoder(w).Encode(data)
		return
	}

	jsonData, err := json.Marshal(Template)
	if err != nil {
		data.Return = 3
		json.NewEncoder(w).Encode(data)
		return
	}

	hash := sha256.Sum256(jsonData)
	data.Hashsum = hex.EncodeToString(hash[:])
	data.Schedule = read.Generator.Generate(Template)

	enc := json.NewEncoder(w)
	enc.SetIndent("", "\t")
	enc.Encode(data)
}
