package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

var (
	port     *int
	interval *int
	urls     []string
	insecure *bool
)

func main() {
	parseFlagsAndArgs()
	registerUI()
	registerAPI()
	listenAndServe()
}

func parseFlagsAndArgs() {
	port = flag.Int("port", 8080, "The port on which to access the results.")
	interval = flag.Int("interval", 60, "The interval between each UI refresh in seconds.")
	insecure = flag.Bool("insecure", false, "If set, will not verify the servers' certificate chain and host name.")
	flag.Parse()
	urls = flag.Args()
}

// registerUI registers an HTTP handler function that presents the results
// in a web user interface
func registerUI() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		results := testNodes(urls)

		tpl, err := template.ParseFiles("index.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		tpl.Execute(w, &ViewModel{*interval, results})
	})
}

// registerAPI registers an HTTP handler function that provides  the results
// in JSON format
func registerAPI() {
	http.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
		results := testNodes(urls)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
	})
}

// testNodes calls the getHTTPStatus method for
// all provided URLs
func testNodes(urls []string) *TestResults {
	results := TestResults{}

	for _, url := range urls {
		result := TestResult{}

		code, err := getHTTPStatus(url)
		if err != nil {
			result.Note = err.Error()
		}
		result.URL = url
		result.Code = code

		if code >= 200 && code < 300 {
			result.IsOK = true
		}
		results = append(results, result)
	}
	return &results
}

// getHTTPStatus issues an HTTP GET call to the specified URL
// and returns the HTTP status code. 0 is returned along with
// an error if the HTTP call could not be completed successfully.
func getHTTPStatus(url string) (code int, err error) {
	res, err := httpClient().Get(url)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()
	return res.StatusCode, nil
}

func httpClient() *http.Client {

	if *insecure {
		transCfg := &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
		return &http.Client{Transport: transCfg}
	}
	return http.DefaultClient
}

func listenAndServe() {
	log.Println("Listening on port " + strconv.Itoa(*port))
	http.ListenAndServe(":"+strconv.Itoa(*port), nil)
}

// TestResult ...
type TestResult struct {
	// The URL under test
	URL string
	// True if the status code is 2xx
	IsOK bool
	// HTTP Status code
	Code int
	// Notes such as error messages
	Note string
}

// TestResults ...
type TestResults []TestResult

// ViewModel is the data structure for the UI template.
type ViewModel struct {
	// Polling interval in seconds
	PollingInterval int
	Results         *TestResults
}
