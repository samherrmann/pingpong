package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/samherrmann/pingpong/config"
	"github.com/samherrmann/pingpong/network"
	"github.com/samherrmann/pingpong/ping"
)

var (
	port           *int
	interval       *int
	insecure       *bool
	configFileName *string

	// timeout is the network node-polling
	// timeout in seconds.
	timeout int

	// version is the version of this software,
	// and should be initialized through the
	// following build flags:
	// go build -ldflags "-X main.version=version"
	version = "latest"
)

func main() {
	log.Println("Running pingpong version " + version)

	parseFlags()
	postFlagsParsingInit()

	nodes, err := parseConfigFile()
	if err != nil {
		log.Printf("Error while parsing config file: %v", err)
		return
	}

	nodesBuff := network.NewNodesBuffer(nodes)
	monitorNodes(nodesBuff)

	err = registerUI(nodesBuff)
	if err != nil {
		log.Printf("Error while registering UI: %v", err)
		return
	}

	registerAPI(nodesBuff)

	err = listenAndServe()
	if err != nil {
		log.Printf("Error while listening for requests: %v", err)
		return
	}
}

// parseConfigFile decodes the configuration JSON
// file into a map.
func parseConfigFile() (*network.Nodes, error) {
	// Attempt to parse config file. If successful,
	// exit immediately.
	nodes, parseErr := config.ParseFile(*configFileName)
	if parseErr == nil {
		return nodes, nil
	}

	// Does file actually exist? If yes, return parse-error
	if _, err := os.Stat(*configFileName); err == nil {
		return nil, parseErr
	}

	// Are we attempting to load the default config file?
	// If no, return parse-error
	if *configFileName != config.DefaultFileName {
		return nil, parseErr
	}

	// Default config file is requested but does
	// not exist, so let's help the user out and
	// create one for them.
	err := config.WriteDefaultFile()
	if err != nil {
		return nil, err
	}

	// Now try parsing the generated config file.
	nodes, err = config.ParseFile(*configFileName)
	if err != nil {
		return nil, err
	}
	log.Println("Config file \"" + *configFileName + "\" was not found. A default file was created and loaded instead.")
	return nodes, nil
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

// monitorNodes calls the getHTTPStatus function for
// all addresses in the config file that have a "http"
// prefix, and calls the ping function for all other
// provided addresses. This function creates separate
// goroutines for each network node under test and is
// therefore non-blocking.
func monitorNodes(nb *network.NodesBuffer) {
	nodes := nb.Get()

	for _, n := range *nodes {
		go func(node network.Node) {
			for {
				var err error
				if node.Method == network.MonitorMethodHTTP {
					_, err = getHTTPStatus(node.URL)
				}
				if node.Method == network.MonitorMethodPing {
					err = ping.Run(node.URL, timeout)
				}
				if err != nil {
					node.Note = err.Error()
					node.IsOK = false
				} else {
					node.Note = ""
					node.IsOK = true
				}
				nb.Update(&node)
				time.Sleep(time.Duration(*interval) * time.Second)
			}
		}(n)
	}
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

func listenAndServe() error {
	log.Println("Listening on port " + strconv.Itoa(*port))
	return http.ListenAndServe(":"+strconv.Itoa(*port), nil)
}
