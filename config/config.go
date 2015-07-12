package config

import (
	"encoding/json"
	"os"
)

const (
	DefaultFileName = "config.json"
)

var (
	Nodes map[string]string
)

// ParseFile decodes the configuration JSON
// file into a map.
func ParseFile(fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&Nodes)
	return err
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
