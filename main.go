package main

import (
	"html/template"
	"net/http"
	"path"
	"log"
	"io/ioutil"
)

type Profile struct {
	Name    string
	Hobbies []string
}

func main() {
	http.HandleFunc("/", foo)
	http.ListenAndServe(":8001", nil)
}

func foo(w http.ResponseWriter, r *http.Request) {
	profile := Profile{"Alex", []string{"snowboarding", "programming", getServerIP()}}


	fp := path.Join("templates", "index.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, profile); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getServerIP() string{
		readApi, err := http.Get("https://api.ipify.org")
		if err != nil {log.Fatal(err)}
		bytes, err := ioutil.ReadAll(readApi.Body)
		if err != nil {log.Fatal(err)}
		return string(bytes)
	}
}
