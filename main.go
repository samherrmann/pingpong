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
