package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	const jsonFilePath = "../blacklist-yt/blacklist.json"

	inputFile, err := os.Open(jsonFilePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer inputFile.Close()

	decoder := json.NewDecoder(inputFile)

	var jsonData map[string]interface{}
	if err := decoder.Decode(&jsonData); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	uniqueKeys := make(map[string]bool)
	for key := range jsonData {
		if uniqueKeys[key] {
			delete(jsonData, key)
		} else {
			uniqueKeys[key] = true
		}
	}

	outputFile, err := os.Create(jsonFilePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer outputFile.Close()

	encoder := json.NewEncoder(outputFile)

	if err := encoder.Encode(jsonData); err != nil {
		fmt.Println("Error encoding and writing JSON:", err)
		return
	}

	fmt.Println("Duplicates removed successfully.")
}