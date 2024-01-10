package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type BlacklistEntry struct {
	// Define a struct with only necessary fields
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

func main() {
	filePath := "../blacklist-yt/blacklist.json"

	// Open the file for reading
	inputFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer inputFile.Close()

	decoder := json.NewDecoder(inputFile)

	// Use a slice with a predefined capacity for seenKeys
	const maxKeys = 100 // Adjust the capacity based on your needs
	var seenKeys = make([]string, 0, maxKeys)
	var data map[string]interface{}
	if err := decoder.Decode(&data); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	// Filtering out duplicate keys
	for key := range data {
		found := false
		for i := 0; i < len(seenKeys); i++ {
			if seenKeys[i] == key {
				delete(data, key)
				found = true
				break
			}
		}
		if !found {
			seenKeys = append(seenKeys, key)
		}
	}

	// Open the file for writing
	outputFile, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer outputFile.Close()

	encoder := json.NewEncoder(outputFile)

	// Encode and write the filtered data
	if err := encoder.Encode(data); err != nil {
		fmt.Println("Error encoding and writing JSON:", err)
		return
	}

	fmt.Println("Duplicates removed successfully.")
}