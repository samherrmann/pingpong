package main

// ViewModel is the data structure for the UI template.
type ViewModel struct {
	// Polling interval in seconds.
	PollingInterval int
	// States is a slice of node states.
	States *NodeStates
	// Version is the version of this software.
	Version string
}
