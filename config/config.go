package config

import (
	"encoding/json"
	"os"
	"sort"
	"strings"

	"github.com/samherrmann/pingpong/network"
)

const (
	DefaultFileName = "config.json"
)

var (
	nodes *map[string]string
)

func Nodes() *network.Nodes {
	states := new(network.Nodes)
	for name, url := range *nodes {
		state := new(network.Node)
		state.Name = name
		state.URL = url

		if strings.HasPrefix(state.URL, "http") {
			state.Method = network.MonitorMethodHTTP
		} else {
			state.Method = network.MonitorMethodPing
		}
		*states = append(*states, *state)
	}
	sort.Sort(states)
	return states
}

// ParseFile decodes the configuration JSON
// file into a map.
func ParseFile(fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&nodes)
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
