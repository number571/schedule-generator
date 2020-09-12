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
	http.HandleFunc("/update", updatePage)
	http.ListenAndServe(":7545", nil)
}

// Method: GET
// Result:
/*
	{
		return: int
	}
*/
// Для теста подключения.
func indexPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct{
		Return int `json:"return"`
	}{
		Return: 0,
	})
}

// Method: POST
// API:
/*
	{
		day: int,      // день недели [0, 6]
		groups: json,  // json с группами
		teachers: json // json с учителями
	}
*/
// Result:
/*
	{
		generator: json, // новый генератор
		return: int      // ошибка
	}
*/
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
// API:
/*
	{
		generator: json // старая версия генератора
	}
*/
// Result:
/*
	{
		generator: json, // новая версия генератора
		schedule: json,  // расписание на день
		hashsum: string, // хеш-сумма шаблона
		return: int      // ошибка
	}
*/
func updatePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var data struct {
		Generator *sc.Generator `json:"generator"`
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

	generator := read.Generator
	hash := sha256.Sum256(jsonData)

	data.Schedule = generator.Generate(Template)
	data.Hashsum = hex.EncodeToString(hash[:])
	data.Generator = generator

	enc := json.NewEncoder(w)
	enc.SetIndent("", "\t")
	enc.Encode(data)
}
