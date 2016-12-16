package httpserver

// ViewModel is the data structure for the UI template.
type ViewModel struct {
	// Polling interval in seconds.
	PollingInterval int
	// Nodes is a slice of newtork nodes.
	Nodes Nodes
	// Version is the version of this software.
	Version string
}

// Node represents a network node.
type Node interface {
	Name() string
	Addr() string
	IsOk() bool
	Note() string
}

// Nodes is a Node slice.
type Nodes []Node
