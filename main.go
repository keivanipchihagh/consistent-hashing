package main

import (
	"fmt"

	"github.com/keivanipchihagh/consistent-hashing/pkg/models"
)

func main() {
	hr := models.NewHashRing(1)

	// Add nodes
	hr.AddServer("S01")
	hr.AddServer("S02")
	hr.AddServer("S03")
	hr.AddServer("S04")

	// Test the distribution of keys
	keys := []string{"C01", "C02", "C03", "C04", "C05", "C06", "C07", "C08", "C09"}

	for _, key := range keys {
		hr.AddClient(key)
	}
	hr.Print()

	fmt.Println("\nRemove S2:")
	hr.RemoveServer("S04")
	hr.DistributeClients()
	hr.Print()
}
