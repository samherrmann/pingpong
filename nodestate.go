package main

import "sync"

// NodeState is a representation of a network node
// (ex: service, server, hub, etc...)
type NodeState struct {
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

// NodeStates is a NodeState slice
type NodeStates []NodeState

func (n NodeStates) Len() int {
	return len(n)
}

func (n NodeStates) Less(i, j int) bool {
	return n[i].Name < n[j].Name
}

func (n NodeStates) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

func NewNodeStatesBuffer() *NodesStatesBuffer {
	return &NodesStatesBuffer{
		mutex: &sync.Mutex{},
	}
}

type NodesStatesBuffer struct {
	states *NodeStates
	mutex  *sync.Mutex
}

func (b *NodesStatesBuffer) Get() *NodeStates {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	return b.states
}

func (b *NodesStatesBuffer) Set(states *NodeStates) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.states = states
}
