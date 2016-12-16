package config

import (
	"github.com/samherrmann/gowrap/jsonfile"
	"github.com/samherrmann/pingpong/network"
)

const (
	// DefaultFileName is the file name
	// of the default configuration file.
	DefaultFileName = "config.json"
)

// ParseFile decodes the configuration JSON
// file into a map.
func ParseFile(fileName string) (*network.Nodes, error) {
	nodes := make(map[string]string)
	err := jsonfile.Read(fileName, &nodes)
	if err != nil {
		return nil, err
	}
	return network.MakeNodes(nodes), nil
}

// WriteDefaultFile creates a file with the default
// configuration.
func WriteDefaultFile() error {
	return jsonfile.Write(DefaultFileName, Default())
}

// Default returns a default configuration.
func Default() map[string]string {
	nodes := make(map[string]string)
	nodes["pingpong host"] = "localhost"
	nodes["pingpong"] = "http://localhost:8080"
	return nodes
}
