package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	f, err := os.Open("./blacklist.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	var d map[string]interface{}
	if err := json.NewDecoder(f).Decode(&d); err != nil {
		fmt.Println(err)
		return
	}

	m := make(map[string]bool)
	for k := range d {
		if m[k] {
			delete(d, k)
		} else {
			m[k] = true
		}
	}

	f, err = os.Create("./blacklist.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(d); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Duplicates removed successfully.")
}
