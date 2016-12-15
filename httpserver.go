package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/samherrmann/pingpong/network"
)

//go:generate go run build/template/main.go

// registerUI registers an HTTP handler function that presents the results
// in a web user interface
func registerUI(nb *network.NodesBuffer) error {
	tpl, err := template.New("index.html").Parse(indexHTML)
	if err != nil {
		return err
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, &ViewModel{*interval, nb.Get(), version})
	})
	return nil
}

// registerAPI registers an HTTP handler function that provides  the results
// in JSON format
func registerAPI(nb *network.NodesBuffer) {
	http.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(nb.Get())
	})
}

func listenAndServe(nb *network.NodesBuffer) error {
	registerAPI(nb)

	err := registerUI(nb)
	if err != nil {
		return err
	}

	log.Println("Listening on port " + strconv.Itoa(*port))
	return http.ListenAndServe(":"+strconv.Itoa(*port), nil)
}
