package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Name struct {
	Name string `json:"name"`
}

func main() {
	log.Println("Starting Name Server")
	http.Handle("/", http.FileServer(http.Dir("./assets")))

	http.HandleFunc("/save", handleSave)
	http.HandleFunc("/names", handleNames)

	log.Println("Serving on http://localhost:8766")
	http.ListenAndServe(":8766", nil)
}

func handleSave(w http.ResponseWriter, r *http.Request) {
	var names []Name
	b, err := os.ReadFile("./data/names.json")
	if err != nil {
		fmt.Println(fmt.Errorf("warn: %v", err))
	} else {
		json.Unmarshal(b, &names)
	}

	err = r.ParseForm()
	if err != nil {
		fmt.Println(fmt.Errorf("error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	names = append(names, Name{
		Name: r.Form.Get("name"),
	})

	bytes, err := json.Marshal(names)
	if err != nil {
		fmt.Println(fmt.Errorf("error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	os.WriteFile("./data/names.json", bytes, os.ModeAppend)
	http.Redirect(w, r, "/", http.StatusFound)
}

func handleNames(w http.ResponseWriter, r *http.Request) {
	b, _ := os.ReadFile("./data/names.json")
	w.Write(b)
}
