package models

import (
	"fmt"
	"sort"
	"sync"

	"github.com/keivanipchihagh/consistent-hashing/internal/utils"
)

type HashRing struct {
	replicas int
	servers  []Server
	clients  []Client
	mapping  map[string]string // Client -> Server
	mu       sync.Mutex
}

func NewHashRing(replicas int) *HashRing {
	return &HashRing{
		replicas: replicas,
		servers:  make([]Server, 0),
		clients:  make([]Client, 0),
		mapping:  make(map[string]string),
		mu:       sync.Mutex{},
	}
}

func (hr *HashRing) Len() int {
	return len(hr.servers)
}

func (hr *HashRing) Less(i, j int) bool {
	return hr.servers[i].hash < hr.servers[j].hash
}

func (hr *HashRing) Swap(i, j int) {
	hr.servers[i], hr.servers[j] = hr.servers[j], hr.servers[i]
}

func (hr *HashRing) AddServer(address string) {
	hr.mu.Lock()
	defer hr.mu.Unlock()

	servers := Server{
		Address:   address,
		hash:      utils.Hash(address),
		isVirtual: false,
	}
	hr.servers = append(hr.servers, servers)
	sort.Sort(hr)
}

func (hr *HashRing) RemoveServer(address string) {
	hr.mu.Lock()
	defer hr.mu.Unlock()

	hash := utils.Hash(address)

	for i, server := range hr.servers {
		if server.hash == hash {
			hr.servers = append(hr.servers[:i], hr.servers[i+1:]...)
			break
		}
	}
	sort.Sort(hr)
}

func (hr *HashRing) GetServer(address string) Server {
	hr.mu.Lock()
	defer hr.mu.Unlock()

	hash := utils.Hash(address)

	idx := sort.Search(len(hr.servers), func(i int) bool {
		return hr.servers[i].hash >= hash
	})

	if idx == len(hr.servers) {
		// Wrap around to the first hash
		idx = 0
	}

	return hr.servers[idx]
}

func (hr *HashRing) AddClient(address string) {

	server := hr.GetServer(address)
	hr.mapping[address] = server.Address
	client := Client{
		Address: address,
		hash:    utils.Hash(address),
	}
	hr.clients = append(hr.clients, client)
}

func (hr *HashRing) DistributeClients() {
	for _, client := range hr.clients {
		server := hr.GetServer(client.Address)
		hr.mapping[client.Address] = server.Address
	}
}

func (hr *HashRing) Print() {
	for _, server := range hr.servers {
		fmt.Printf("%s: ", server.Address)
		for clientAddr, serverAddr := range hr.mapping {
			if serverAddr == server.Address {
				fmt.Printf("%s ", clientAddr)
			}
		}
		fmt.Println()
	}
}
