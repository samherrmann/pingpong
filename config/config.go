package config

import (
	"encoding/json"
	"os"

	"github.com/samherrmann/pingpong/network"
)

const (
	DefaultFileName = "config.json"
)

// ParseFile decodes the configuration JSON
// file into a map.
func ParseFile(fileName string) (*network.Nodes, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	nodes := make(map[string]string)
	err = json.NewDecoder(file).Decode(&nodes)
	if err != nil {
		return nil, err
	}
	return network.MakeNodes(nodes), nil
}

func WriteDefaultFile() error {
	file, err := os.Create(DefaultFileName)
	if err != nil {
		return err
	}
	defer file.Close()

	json, err := json.MarshalIndent(defaultConfig(), "", "    ")
	if err != nil {
		return err
	}
	_, err = file.Write(json)
	return err
}

func defaultConfig() map[string]string {
	nodes := make(map[string]string)
	nodes["pingpong host"] = "localhost"
	nodes["pingpong"] = "http://localhost:8080"
	return nodes
}
