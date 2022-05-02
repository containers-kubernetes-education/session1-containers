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
	b, err := os.ReadFile("./names.json")
	if err != nil {
		fmt.Println(fmt.Errorf("error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = r.ParseForm()
	if err != nil {
		fmt.Println(fmt.Errorf("error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.Unmarshal(b, &names)
	names = append(names, Name{
		Name: r.Form.Get("name"),
	})

	bytes, err := json.Marshal(names)
	if err != nil {
		fmt.Println(fmt.Errorf("error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	os.WriteFile("./names.json", bytes, os.ModeAppend)
	http.Redirect(w, r, "/", http.StatusFound)
}

func handleNames(w http.ResponseWriter, r *http.Request) {
	b, err := os.ReadFile("./names.json")
	if err != nil {
		fmt.Println(fmt.Errorf("error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(b)
}
