package main

import (
	"fmt"

	"github.com/keivanipchihagh/consistent-hashing/pkg/models"
)

func main() {
	hr := models.NewHashRing()

	// Add nodes
	hr.AddNode("Node1")
	hr.AddNode("Node2")
	hr.AddNode("Node3")

	// Test the distribution of keys
	keys := []string{"key1", "key2", "key3", "key4", "key5"}

	for _, key := range keys {
		node := hr.GetNode(key)
		fmt.Printf("Key %s is mapped to node %s\n", key, node.Address)
	}

	// Remove a node and see how keys are reassigned
	fmt.Println("\nRemoving Node2")
	hr.RemoveNode("Node2")

	for _, key := range keys {
		node := hr.GetNode(key)
		fmt.Printf("Key %s is now mapped to node %s\n", key, node.Address)
	}
}
