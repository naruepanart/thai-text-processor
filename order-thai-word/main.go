package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	// Read the blacklist JSON file into a map of strings to interfaces
	d, err := readJSONFile("./blacklist.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Remove any duplicate entries from the map
	d = removeDuplicates(d)

	// Write the updated blacklist JSON file
	err = writeJSONFile("./blacklist.json", d)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func readJSONFile(path string) (map[string]interface{}, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var d map[string]interface{}
	if err := json.NewDecoder(f).Decode(&d); err != nil {
		return nil, err
	}

	return d, nil
}

func writeJSONFile(path string, d map[string]interface{}) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	err = json.NewEncoder(f).Encode(d)
	if err != nil {
		return err
	}

	return nil
}

func removeDuplicates(d map[string]interface{}) map[string]interface{} {
	m := make(map[string]bool)
	for k := range d {
		if m[k] {
			delete(d, k)
		} else {
			m[k] = true
		}
	}
	return d
}
