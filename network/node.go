package network

import (
	"sort"
	"strings"
	"sync"
)

const (
	MonitorMethodHTTP = "HTTP/S"
	MonitorMethodPing = "Ping"
)

// Node is a representation of a network node
// (ex: service, server, hub, etc...)
type Node struct {
	// The name of the node
	Name string
	// The URL of the network node
	URL string
	// True if the node is operational
	IsOK bool
	// Notes such as error messages
	Note string
	// The method that was used to test
	// the node's availability
	// (ex: ping or http)
	Method string
}

// Nodes is a Node slice
type Nodes []Node

func (n Nodes) Len() int {
	return len(n)
}

func (n Nodes) Less(i, j int) bool {
	return n[i].Name < n[j].Name
}

func (n Nodes) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

// MakeNodes creates a Nodes slice from the provided
// map in which the keys are expected to be the names
// of the nodes and the values are expected to be the
// addresses of the nodes.
func MakeNodes(nodesMap map[string]string) *Nodes {
	ns := new(Nodes)
	for name, url := range nodesMap {
		n := new(Node)
		n.Name = name
		n.URL = url

		if strings.HasPrefix(n.URL, "http") {
			n.Method = MonitorMethodHTTP
		} else {
			n.Method = MonitorMethodPing
		}
		*ns = append(*ns, *n)
	}
	sort.Sort(ns)
	return ns
}

// NewNodesBuffer returns a NodesBuffer struct
func NewNodesBuffer(nodes *Nodes) *NodesBuffer {
	return &NodesBuffer{
		mutex: &sync.Mutex{},
		nodes: nodes,
	}
}

// NodesBuffer allows for safe
// asyn-operations on a Nodes
// slice.
type NodesBuffer struct {
	nodes *Nodes
	mutex *sync.Mutex
}

func (b *NodesBuffer) Get() *Nodes {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	return b.nodes
}

func (b *NodesBuffer) Update(node *Node) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	for i := range *b.nodes {
		if (*b.nodes)[i].URL == node.URL {
			(*b.nodes)[i] = *node
		}
	}
}
