package models

import (
	"hash/crc32"
	"sort"
	"sync"
)

type HashRing struct {
	nodes []Node
	mu    sync.Mutex
}

func NewHashRing() *HashRing {
	return &HashRing{
		nodes: make([]Node, 0),
		mu:    sync.Mutex{},
	}
}

func (hr *HashRing) Len() int {
	return len(hr.nodes)
}

func (hr *HashRing) Less(i, j int) bool {
	return hr.nodes[i].hash < hr.nodes[j].hash
}

func (hr *HashRing) Swap(i, j int) {
	hr.nodes[i], hr.nodes[j] = hr.nodes[j], hr.nodes[i]
}

func (hr *HashRing) AddNode(address string) {
	hr.mu.Lock()
	defer hr.mu.Unlock()

	hash := crc32.ChecksumIEEE([]byte(address))

	node := Node{
		Address: address,
		hash:    hash,
	}

	hr.nodes = append(hr.nodes, node)
	sort.Sort(hr)
}

func (hr *HashRing) RemoveNode(address string) {
	hr.mu.Lock()
	defer hr.mu.Unlock()

	hash := crc32.ChecksumIEEE([]byte(address))

	for i, node := range hr.nodes {
		if node.hash == hash {
			hr.nodes = append(hr.nodes[:i], hr.nodes[i+1:]...)
			break
		}
	}
	sort.Sort(hr)
}

func (hr *HashRing) GetNode(key string) Node {
	hr.mu.Lock()
	defer hr.mu.Unlock()

	if len(hr.nodes) == 0 {
		return Node{}
	}

	hash := crc32.ChecksumIEEE([]byte(key))

	idx := sort.Search(len(hr.nodes), func(i int) bool {
		return hr.nodes[i].hash >= hash
	})

	// If the hash is greater than all, wrap around to the first hash
	if idx == len(hr.nodes) {
		idx = 0
	}

	return hr.nodes[idx]
}
