package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func dedupJSON(input []byte) ([]byte, error) {
	var d map[string]interface{}
	if err := json.Unmarshal(input, &d); err != nil {
		return nil, err
	}

	s := make(map[string]bool)
	for k := range d {
		if s[k] {
			delete(d, k)
		} else {
			s[k] = true
		}
	}

	r, err := json.Marshal(d)
	return r, err
}

func main() {
	f := "your_file.json"
	if i, err := os.ReadFile(f); err != nil {
		fmt.Println("Error reading file:", err)
	} else if o, err := dedupJSON(i); err != nil {
		fmt.Println("Error removing duplicates:", err)
	} else if err := os.WriteFile(f, o, 0644); err != nil {
		fmt.Println("Error writing to file:", err)
	} else {
		fmt.Println("Duplicates removed successfully.")
	}
}