package monitor

import (
	"time"

	"github.com/samherrmann/pingpong/network"
)

// Node represents a network node.
type Node interface {
	Addr() string
	SetIsOk(v bool)
	SetNote(v string)
}

// Nodes is a slice of nodes.
type Nodes []Node

// Run calls the getHTTPStatus function for
// all addresses in the config file that have a "http"
// prefix, and calls the ping function for all other
// provided addresses. This function creates separate
// goroutines for each network node under test and is
// therefore non-blocking.
func Run(interval int, ns Nodes) {
	for _, n := range ns {
		go func(n Node) {
			for {
				execTest(interval, n)
				time.Sleep(time.Duration(interval) * time.Second)
			}
		}(n)
	}
}

func execTest(timeout int, n Node) {
	err := network.Ping(n.Addr(), timeout)
	if err != nil {
		n.SetIsOk(false)
		n.SetNote(err.Error())
	} else {
		n.SetIsOk(true)
		n.SetNote("")
	}
}

// // getHTTPStatus issues an HTTP GET call to the specified URL
// // and returns the HTTP status code. 0 is returned along with
// // an error if the HTTP call could not be completed successfully.
// func getHTTPStatus(url string) (code int, err error) {
// 	res, err := httpClient().Get(url)
// 	if err != nil {
// 		return 0, err
// 	}
// 	defer res.Body.Close()
// 	return res.StatusCode, nil
// }

// // httpClient returns a HTTP client that does not
// // verify a server's certificate chain and host name
// // if the "insecure" command-line flag in set to true.
// func httpClient() *http.Client {

// 	timeout := time.Duration(timeout) * time.Second

// 	if *insecure {
// 		transport := &http.Transport{
// 			TLSClientConfig: &tls.Config{
// 				InsecureSkipVerify: true,
// 			},
// 		}
// 		return &http.Client{
// 			Transport: transport,
// 			Timeout:   timeout,
// 		}
// 	}
// 	return &http.Client{Timeout: timeout}
// }
