package main

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/samherrmann/pingpong/network"
	"github.com/samherrmann/pingpong/ping"
)

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
