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
	bl := loadBlacklist()

	// Compile the blacklist into regexps and replacement strings
	rs, rb := compileBlacklist(bl)

	// Get a list of all text files in the current directory
	tf := getTextFiles()

	// Loop through the text files
	for _, name := range tf {
		// Read the file into a byte slice
		data := readFile(name)

		// Update the content of the byte slice using the regexps and replacement strings
		data = updateContent(data, rs, rb)

		// Convert any years found in the text to the Buddhist Era
		data = convertYears(data)

		// Write the updated byte slice back to the file
		writeFile(name, data)
	}
}

func loadBlacklist() map[string]string {
	f, err := os.Open("blacklist.json")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer f.Close()

	bl := make(map[string]string)
	err = json.NewDecoder(f).Decode(&bl)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return bl
}

func compileBlacklist(bl map[string]string) ([]*regexp.Regexp, [][]byte) {
	var rs []*regexp.Regexp
	var rb [][]byte
	for p, r := range bl {
		re, err := regexp.Compile(p)
		if err != nil {
			fmt.Println(err)
			return nil, nil
		}
		rs = append(rs, re)
		rb = append(rb, []byte(r))
	}
	return rs, rb
}

func getTextFiles() []string {
	tf, err := filepath.Glob("*.txt")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return tf
}

func readFile(name string) []byte {
	data, err := os.ReadFile(name)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return data
}

func updateContent(data []byte, rs []*regexp.Regexp, rb [][]byte) []byte {
	for i, r := range rs {
		data = r.ReplaceAll(data, rb[i])
	}
	return data
}

func convertYears(data []byte) []byte {
	r := regexp.MustCompile(`(\d{4})`)
	return r.ReplaceAllFunc(data, func(m []byte) []byte {
		y, _ := strconv.Atoi(string(m))
		if y >= 2400 {
			return []byte(strconv.Itoa(y - 543))
		}
		return m
	})
}

func writeFile(name string, data []byte) {
	err := os.WriteFile(name, data, os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}
}