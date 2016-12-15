package main

import "github.com/samherrmann/pingpong/network"

// ViewModel is the data structure for the UI template.
type ViewModel struct {
	// Polling interval in seconds.
	PollingInterval int
	// Nodes is a slice of newtork nodes.
	Nodes *network.Nodes
	// Version is the version of this software.
	Version string
}
