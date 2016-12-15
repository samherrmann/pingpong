package main

import (
	"flag"

	"github.com/samherrmann/pingpong/config"
)

var (
	port           *int
	interval       *int
	insecure       *bool
	configFileName *string

	// timeout is the network node-polling
	// timeout in seconds.
	timeout int
)

func parseFlags() {
	port = flag.Int("port", 8080, "The port on which to access the results.")
	interval = flag.Int("interval", 60, "The interval between each UI refresh in seconds.")
	insecure = flag.Bool("insecure", false, "If set, will not verify the servers' certificate chain and host name.")
	configFileName = flag.String("config", config.DefaultFileName, "The name of the configuration file.")
	flag.Parse()
}

// postFlagsParsingInit performs initialization work
// that depends on values that are provided through
// command-line flags.
func postFlagsParsingInit() {
	timeout = *interval / 2
}
