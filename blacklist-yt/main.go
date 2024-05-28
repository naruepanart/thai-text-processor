package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

func main() {
	// Load the blacklist from the "blacklist.json" file
	bl := make(map[string]string)
	f, err := os.Open("blacklist.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	err = json.NewDecoder(f).Decode(&bl)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Compile the blacklist into regexps and replacement strings
	var rs []*regexp.Regexp
	var rb [][]byte
	for p, r := range bl {
		// Compile the pattern string into a regex
		re, err := regexp.Compile(p)
		if err != nil {
			fmt.Println(err)
			return
		}
		// Append the regex to the slice of regexps
		rs = append(rs, re)
		// Append the replacement string as a byte slice to the slice of byte slices
		rb = append(rb, []byte(r))
	}

	// Get a list of all text files in the current directory
	tf, err := filepath.Glob("*.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Loop through the text files
	for _, name := range tf {
		// Read the file into a byte slice
		data, err := os.ReadFile(name)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Update the content of the byte slice using the regexps and
		// replacement strings
		for i, r := range rs {
			// Replace all occurrences of the regex in the byte slice with the
			// corresponding replacement string from the slice of byte slices
			data = r.ReplaceAll(data, rb[i])
		}

		// Convert any years found in the text to the Buddhist Era
		r := regexp.MustCompile(`(\d{4})`)
		data = r.ReplaceAllFunc(data, func(m []byte) []byte {
			y, _ := strconv.Atoi(string(m))
			if y >= 2400 {
				return []byte(strconv.Itoa(y - 543))
			}
			return m
		})

		// Write the updated byte slice back to the file
		err = os.WriteFile(name, data, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
