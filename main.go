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
	"time"
)

var (
	port             *int
	interval         *int
	insecure         *bool
	configFileName   *string
	config           map[string]string
	nodeStatesBuffer = NewNodeStatesBuffer()
)

func main() {
	parseFlagsAndArgs()
	parseConfigFile()
	registerUI()
	registerAPI()
	monitorNodes()
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
	tpl, err := template.ParseFiles("index.html")
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

			for name, url := range config {
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

					err := ping(state.URL)
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

// httpClient returns a HTTP client that does not
// verify a server's certificate chain and host name
// if the "insecure" command-line flag in set to true.
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
