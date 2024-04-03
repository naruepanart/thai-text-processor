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

	var data map[string]interface{}
	if err := decoder.Decode(&data); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	// Filtering out duplicate keys using a hash table
	seenKeys := make(map[string]bool)
	filteredData := make(map[string]interface{})
	for key, value := range data {
		if !seenKeys[key] {
			filteredData[key] = value
			seenKeys[key] = true
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
	if err := encoder.Encode(filteredData); err != nil {
		fmt.Println("Error encoding and writing JSON:", err)
		return
	}

	fmt.Println("Duplicates removed successfully.")
}