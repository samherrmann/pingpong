package main

import (
	"log"

	"github.com/samherrmann/pingpong/app/about"
	"github.com/samherrmann/pingpong/app/config"
	"github.com/samherrmann/pingpong/app/httpserver"
	"github.com/samherrmann/pingpong/app/monitor"
	"github.com/samherrmann/pingpong/app/node"
)

//go:generate go run build/template/main.go

func main() {
	a := about.New(version)
	a.Log()

	c, err := config.ImportOrExportSample()
	if err != nil {
		log.Printf("Error while loading config file: %v", err)
		return
	}

	nodes := node.Slice(c.Nodes)
	monitorNodes(c.Interval, nodes)

	err = runServer(c.Port, c.Interval, nodes)
	if err != nil {
		log.Printf("Error while setting up web server: %v", err)
		return
	}
}

func monitorNodes(interval int, n node.Nodes) {
	// convert []*node.Nodes to []*monitor.Nodes
	nodes := make(monitor.Nodes, len(n))
	for i, v := range n {
		nodes[i] = v
	}
	monitor.Run(interval, nodes)
}

func runServer(port int, interval int, n node.Nodes) error {
	// convert []*node.Nodes to []*httpserver.Nodes
	nodes := make(httpserver.Nodes, len(n))
	for i, v := range n {
		nodes[i] = v
	}
	return httpserver.ListenAndServe(port, interval, version, nodes)
}
