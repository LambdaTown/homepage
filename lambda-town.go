package main

import (
	_ "embed"
	"encoding/json"
	"html/template"
	"net/http"
	"time"
)

//go:embed index.html
var INDEX string

//go:embed metadata.json
var METADATA string

type Metadata struct {
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Events      []time.Time `json:"events"`
}

type TemplateInfill struct {
	Title        string
	Description  string
	PastEvents   []time.Time
	FutureEvents []time.Time
}

func partitionTimes(threshold time.Time, times []time.Time) (past []time.Time, future []time.Time) {
	for _, t := range times {
		if threshold.Before(t) {
			future = append(future, t)
		} else {
			past = append(past, t)
		}
	}
	return
}

func home(w http.ResponseWriter, r *http.Request) {
	template, err := template.New("home template").Parse(INDEX)
	if err != nil {
		panic(err)
	}
	metadata := Metadata{}
	err = json.Unmarshal([]byte(METADATA), &metadata)
	if err != nil {
		panic(err)
	}
	past, future := partitionTimes(time.Now(), metadata.Events)
	infill := TemplateInfill{metadata.Title, metadata.Description, past, future}
	err = template.Execute(w, infill)
	if err != nil {
		panic(err)
	}
}

func api(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(METADATA))
}

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/metadata.json", api)
	http.ListenAndServe(":2137", nil)
}
