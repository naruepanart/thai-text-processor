package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	filePath := "../blacklist-yt/blacklist.json"

	input, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	var data map[string]interface{}
	if err := json.Unmarshal(input, &data); err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	// Use a stack-allocated array for seenKeys
	seenKeys := make(map[string]bool, len(data))
	for key := range data {
		if seenKeys[key] {
			delete(data, key)
		} else {
			seenKeys[key] = true
		}
	}

	output, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	err = os.WriteFile(filePath, output, 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("Duplicates removed successfully.")
}