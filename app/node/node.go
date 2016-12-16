package node

import "sort"
import "sync"

// New returns a Node struct.
func New(name string, addr string) *Node {
	n := new(Node)
	n.mutex = new(sync.RWMutex)
	n.SetName(name)
	n.SetAddr(addr)
	return n
}

// data contains the Node data.
type data struct {
	Name string `json:"name"`
	Addr string `json:"addr"`
	IsOK bool   `json:"isOK"`
	Note string `json:"note"`
}

// Node is a representation of a network node.
// Ex: service, server, hub, etc...
type Node struct {
	mutex *sync.RWMutex
	data
}

// Name returns the name of the node.
func (n *Node) Name() string {
	n.mutex.RLock()
	defer n.mutex.RUnlock()
	return n.data.Name
}

// SetName sets the name of the node.
func (n *Node) SetName(v string) {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	n.data.Name = v
}

// Addr returns the address of the node.
func (n *Node) Addr() string {
	n.mutex.RLock()
	defer n.mutex.RUnlock()
	return n.data.Addr
}

// SetAddr sets the address of the node.
func (n *Node) SetAddr(v string) {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	n.data.Addr = v
}

// IsOk returns true if the
// node is operational.
func (n *Node) IsOk() bool {
	n.mutex.RLock()
	defer n.mutex.RUnlock()
	return n.data.IsOK
}

// SetIsOk sets the OK state of
// the note.
func (n *Node) SetIsOk(v bool) {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	n.data.IsOK = v
}

// Note returns notes such as error
// messages  about the node
func (n *Node) Note() string {
	n.mutex.RLock()
	defer n.mutex.RUnlock()
	return n.data.Note
}

// SetNote sets a note for the node.
func (n *Node) SetNote(v string) {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	n.data.Note = v
}

// Nodes is a collection of network nodes.
type Nodes []*Node

// Slice returns a Nodes slice.
// The slice is alpha-numerically sorted
// by the the node name.
func Slice(m map[string]string) Nodes {
	ns := make(Nodes, 0)

	for name, addr := range m {
		n := New(name, addr)
		ns = append(ns, n)
	}
	sort.Sort(ns)
	return ns
}

func (n Nodes) Len() int {
	return len(n)
}

func (n Nodes) Less(i, j int) bool {
	return n[i].Name() < n[j].Name()
}

func (n Nodes) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}
