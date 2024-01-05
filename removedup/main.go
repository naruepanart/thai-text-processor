package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func dedupJSON(input []byte) ([]byte, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(input, &data); err != nil {
		return nil, err
	}

	seenKeys := make(map[string]bool)
	for key := range data {
		if seenKeys[key] {
			delete(data, key)
		} else {
			seenKeys[key] = true
		}
	}

	result, err := json.Marshal(data)
	return result, err
}

func main() {
	filePath := "../blacklist-yt/blacklist.json"
	if input, err := os.ReadFile(filePath); err != nil {
		fmt.Println("Error reading file:", err)
	} else if output, err := dedupJSON(input); err != nil {
		fmt.Println("Error removing duplicates:", err)
	} else if err := os.WriteFile(filePath, output, 0644); err != nil {
		fmt.Println("Error writing to file:", err)
	} else {
		fmt.Println("Duplicates removed successfully.")
	}
}