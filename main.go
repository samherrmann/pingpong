package main

import (
	"log"

	"github.com/samherrmann/pingpong/network"
)

func main() {
	log.Println("Running pingpong version " + version)
	parseFlags()
	nodes, err := parseConfigFile()
	if err != nil {
		log.Printf("Error while parsing config file: %v", err)
		return
	}
	nodesBuff := network.NewNodesBuffer(nodes)
	monitorNodes(nodesBuff)

	err = listenAndServe(nodesBuff)
	if err != nil {
		log.Printf("Error while setting up web server: %v", err)
		return
	}
}
