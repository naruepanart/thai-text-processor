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
	// This involves opening the file and decoding the JSON data into a
	// map of strings to strings.
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
	// This involves compiling each pattern in the blacklist into a
	// regex and creating a slice of regexps and a slice of byte slices
	// containing the replacement strings.
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

	// Find all the text files in the directory
	// This involves using the filepath.Glob function to find all files
	// in the current directory with the extension ".txt".
	tfs, err := filepath.Glob("*.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Process each text file with the compiled blacklist
	// This involves reading each file into a byte slice, updating the
	// content of the byte slice using the regexps and replacement strings,
	// and then writing the updated byte slice back to the file.
	for _, name := range tfs {
		// Read the file into a byte slice
		data, err := os.ReadFile(name)
		if err != nil {
			// If there's an error, skip this file
			continue
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
		if err := os.WriteFile(name, data, os.ModePerm); err != nil {
			// If there's an error, skip this file
			continue
		}
	}
}