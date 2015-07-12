package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

var (
	port           *int
	interval       *int
	insecure       *bool
	configFileName *string
	config         map[string]string
)

func main() {
	parseFlagsAndArgs()
	parseConfigFile()
	registerUI()
	registerAPI()
	listenAndServe()
}

// parseConfigFile decodes the configuration JSON
// file into a map.
func parseConfigFile() {
	configFile, err := os.Open(*configFileName)
	if err != nil {
		panic(err.Error())
	}

	err = json.NewDecoder(configFile).Decode(&config)
	if err != nil {
		panic(err.Error())
	}
}

func parseFlagsAndArgs() {
	port = flag.Int("port", 8080, "The port on which to access the results.")
	interval = flag.Int("interval", 60, "The interval between each UI refresh in seconds.")
	insecure = flag.Bool("insecure", false, "If set, will not verify the servers' certificate chain and host name.")
	configFileName = flag.String("config", "config.json", "The name of the configuration file.")
	flag.Parse()
}

// registerUI registers an HTTP handler function that presents the results
// in a web user interface
func registerUI() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		results := testNodes()

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
		results := testNodes()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
	})
}

// testNodes calls the getHTTPStatus function for
// all addresses in the config file that have a "http"
// prefix, and calls the ping function for all other
// provided addresses.
func testNodes() *TestResults {
	results := TestResults{}

	for name, url := range config {
		result := TestResult{}
		result.Name = name
		result.URL = url

		if strings.HasPrefix(result.URL, "http") {
			result.Method = "HTTP/S"

			code, err := getHTTPStatus(result.URL)
			if err != nil {
				result.Note = err.Error()
			}

			if code >= 200 && code < 300 {
				result.IsOK = true
			}

		} else {
			result.Method = "Ping"

			err := ping(result.URL)
			if err != nil {
				result.Note = err.Error()
			} else {
				result.IsOK = true
			}
		}

		results = append(results, result)
	}
	sort.Sort(&results)
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

// ping returns nil if the specified IP address
// is responsive to pings, and returns an error
// otherwise
//
// NOTE: Currently only supported on Linux
func ping(ipAddr string) (err error) {
	cmd := exec.Command("ping", "-c", "1", "-w", "1", ipAddr)
	if err := cmd.Start(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
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
	// The name of the service
	Name string
	// The URL under test
	URL string
	// True if the status code is 2xx
	IsOK bool
	// Notes such as error messages
	Note string
	// The method that was used to test
	// the service's availability
	// (ex: ping or http)
	Method string
}

// TestResults ...
type TestResults []TestResult

func (results TestResults) Len() int {
	return len(results)
}

func (results TestResults) Less(i, j int) bool {
	return results[i].Name < results[j].Name
}

func (results TestResults) Swap(i, j int) {
	results[i], results[j] = results[j], results[i]
}

// ViewModel is the data structure for the UI template.
type ViewModel struct {
	// Polling interval in seconds
	PollingInterval int
	Results         *TestResults
}
