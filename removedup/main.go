package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Entry struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

func main() {
	file := "../blacklist-yt/blacklist.json"

	input, err := os.Open(file)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer input.Close()

	decoder := json.NewDecoder(input)

	var data map[string]interface{}
	if err := decoder.Decode(&data); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	seen := make(map[string]bool)
	for key := range data {
		if seen[key] {
			delete(data, key)
		} else {
			seen[key] = true
		}
	}

	output, err := os.Create(file)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer output.Close()

	encoder := json.NewEncoder(output)

	if err := encoder.Encode(data); err != nil {
		fmt.Println("Error encoding and writing JSON:", err)
		return
	}

	fmt.Println("Duplicates removed successfully.")
}