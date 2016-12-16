package config

import (
	"errors"
	"os"

	"github.com/samherrmann/pingpong/jsonfile"
)

const (
	fileName = "pingpong.config.json"
)

// Config is the application configuration
type Config struct {
	Port     int
	Interval int
	Insecure bool
	Nodes    map[string]string
}

// New returns a Config struct.
func New() *Config {
	c := new(Config)
	c.Port = 8080
	c.Interval = 60
	c.Insecure = false
	c.Nodes = make(map[string]string)
	return c
}

// Sample returns a Config with sample nodes.
func Sample() *Config {
	c := New()
	c.Nodes["pingpong host"] = "localhost"
	c.Nodes["pingpong"] = "http://localhost:8080"
	return c
}

// Import attempts to parse the config file
// and load the data into the Config struct.
func (c *Config) Import() error {
	return jsonfile.Read(fileName, c)
}

// Export attempts to write the Config struct
// to the config file.
func (c *Config) Export() error {
	return jsonfile.Write(fileName, c)
}

// DoesFileExist returns true if the config
// file already exists.
func DoesFileExist() bool {
	_, err := os.Stat(fileName)
	return err == nil
}

// FileName returns the name of the
// config file.
func FileName() string {
	return fileName
}

// ImportOrExportSample attempts to import the
// configuration file. If the file does not exits,
// a sample configuration file is exported before
// the function returns an error.
func ImportOrExportSample() (*Config, error) {
	c := New()

	// Attempt to parse config file. If successful,
	// exit immediately.
	parseErr := c.Import()
	if parseErr == nil {
		return c, nil
	}

	// Does file actually exist? If yes, return parse-error
	if DoesFileExist() {
		return nil, parseErr
	}

	// Config file does not exist, so let's help the user
	// out and create one for them.
	c = Sample()
	err := c.Export()
	if err != nil {
		return nil, err
	}

	return nil, errors.New("No '" + FileName() + "' found. " +
		"A sample file was created in the current directory. " +
		"Edit the file as required and re-run pingpong.")
}
