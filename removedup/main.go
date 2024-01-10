package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	filePath := "../blacklist-yt/blacklist.json"

	// Open the file for reading
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	// Use a stack-allocated array for seenKeys
	seenKeys := make(map[string]bool)
	var data map[string]interface{}
	if err := decoder.Decode(&data); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	// Filtering out duplicate keys
	for key := range data {
		if seenKeys[key] {
			delete(data, key)
		} else {
			seenKeys[key] = true
		}
	}

	// Open the file for writing
	file, err = os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)

	// Encode and write the filtered data
	if err := encoder.Encode(data); err != nil {
		fmt.Println("Error encoding and writing JSON:", err)
		return
	}

	fmt.Println("Duplicates removed successfully.")
}
