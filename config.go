package main

import (
	"log"
	"os"

	"github.com/samherrmann/pingpong/config"
	"github.com/samherrmann/pingpong/network"
)

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
