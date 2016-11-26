package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"pingpong/config"
	"pingpong/ping"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	port             *int
	interval         *int
	insecure         *bool
	configFileName   *string
	nodeStatesBuffer = NewNodeStatesBuffer()

	// timeout is the network node-polling
	// timeout in seconds.
	timeout int
)

func main() {
	parseFlags()
	postFlagsParsingInit()
	parseConfigFile()
	registerUI()
	registerAPI()
	monitorNodes()
	listenAndServe()
}

// parseConfigFile decodes the configuration JSON
// file into a map.
func parseConfigFile() {
	// Attempt to parse config file. If successful,
	// exit immediately.
	parseErr := config.ParseFile(*configFileName)
	if parseErr == nil {
		return
	}

	// Does file actually exist?
	// If yes, panic
	if _, err := os.Stat(*configFileName); err == nil {
		panic(parseErr)
	}

	// Are we attempting to load the default config file?
	// If no, panic
	if *configFileName != config.DefaultFileName {
		panic(parseErr)
	}

	// Default config file is requested but does
	// not exist, so let's help the user out and
	// create one for them.
	err := config.WriteDefaultFile()
	if err != nil {
		panic(err)
	}

	// Now try parsing the generated config file.
	err = config.ParseFile(*configFileName)
	if err != nil {
		panic(err)
	}
	log.Println("Config file \"" + *configFileName + "\" was not found. A default file was created and loaded instead.")
}

func parseFlags() {
	port = flag.Int("port", 8080, "The port on which to access the results.")
	interval = flag.Int("interval", 60, "The interval between each UI refresh in seconds.")
	insecure = flag.Bool("insecure", false, "If set, will not verify the servers' certificate chain and host name.")
	configFileName = flag.String("config", config.DefaultFileName, "The name of the configuration file.")
	flag.Parse()
}

// postFlagsParsingInit performs initialization work
// that depends on values that are provided through
// command-line flags.
func postFlagsParsingInit() {
	timeout = *interval / 2
}

//go:generate go run scripts/template.go

// registerUI registers an HTTP handler function that presents the results
// in a web user interface
func registerUI() {
	tpl, err := template.New("index.html").Parse(indexHTML)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, &ViewModel{*interval, nodeStatesBuffer.Get()})
	})
}

// registerAPI registers an HTTP handler function that provides  the results
// in JSON format
func registerAPI() {
	http.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(nodeStatesBuffer.Get())
	})
}

// monitorNodes calls the getHTTPStatus function for
// all addresses in the config file that have a "http"
// prefix, and calls the ping function for all other
// provided addresses. This function is executed in a
// goroutine and is therefore non-blocking.
func monitorNodes() {
	go func() {
		for {
			states := &NodeStates{}

			for name, url := range config.Nodes {
				state := &NodeState{}
				state.Name = name
				state.URL = url

				if strings.HasPrefix(state.URL, "http") {
					state.Method = "HTTP/S"

					code, err := getHTTPStatus(state.URL)
					if err != nil {
						state.Note = err.Error()
					}

					if code >= 200 && code < 300 {
						state.IsOK = true
					}

				} else {
					state.Method = "Ping"

					err := ping.Run(state.URL, timeout)
					if err != nil {
						state.Note = err.Error()
					} else {
						state.IsOK = true
					}
				}

				*states = append(*states, *state)
			}
			sort.Sort(states)
			nodeStatesBuffer.Set(states)
			time.Sleep(time.Duration(*interval) * time.Second)
		}
	}()
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

// httpClient returns a HTTP client that does not
// verify a server's certificate chain and host name
// if the "insecure" command-line flag in set to true.
func httpClient() *http.Client {

	timeout := time.Duration(timeout) * time.Second

	if *insecure {
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
		return &http.Client{
			Transport: transport,
			Timeout:   timeout,
		}
	}
	return &http.Client{Timeout: timeout}
}

func listenAndServe() {
	log.Println("Listening on port " + strconv.Itoa(*port))
	http.ListenAndServe(":"+strconv.Itoa(*port), nil)
}
