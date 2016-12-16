package jsonfile

import (
	"encoding/json"
	"os"
)

// Read decodes a JSON file into a Go value
// pointed to by v.
func Read(filePath string, v interface{}) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewDecoder(file).Decode(v)
}

// Write marshals the values pointed to by v into
// a JSON file.
func Write(filePath string, v interface{}) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	json, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		return err
	}
	_, err = file.Write(json)
	return err
}
