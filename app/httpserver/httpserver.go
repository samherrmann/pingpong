package httpserver

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// registerUI registers an HTTP handler function that presents the results
// in a web user interface
func registerUI(interval int, version string, nodes Nodes) error {
	tpl, err := template.New("index.html").Parse(indexHTML)
	if err != nil {
		return err
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, &ViewModel{interval, nodes, version})
	})
	return nil
}

// registerAPI registers an HTTP handler function that provides  the results
// in JSON format
func registerAPI(nodes Nodes) {
	http.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(nodes)
	})
}

// ListenAndServe starts the HTTP server.
func ListenAndServe(port int, interval int, version string, nodes Nodes) error {
	registerAPI(nodes)
	err := registerUI(interval, version, nodes)
	if err != nil {
		return err
	}

	log.Println("Listening on port " + strconv.Itoa(port))
	return http.ListenAndServe(":"+strconv.Itoa(port), nil)
}
