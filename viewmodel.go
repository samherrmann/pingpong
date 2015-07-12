package main

// ViewModel is the data structure for the UI template.
type ViewModel struct {
	// Polling interval in seconds
	PollingInterval int
	States          *NodeStates
}
