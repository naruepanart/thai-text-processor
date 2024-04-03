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

	max := 100
	var seen = make([]string, 0, max)
	var data map[string]interface{}
	if err := decoder.Decode(&data); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	for key := range data {
		found := false
		for i := 0; i < len(seen); i++ {
			if seen[i] == key {
				delete(data, key)
				found = true
				break
			}
		}
		if !found {
			seen = append(seen, key)
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
