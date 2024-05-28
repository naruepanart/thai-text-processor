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
	bl, err := loadBlacklist()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Compile the blacklist into regexps and replacement strings
	// This involves compiling each pattern in the blacklist into a
	// regex and creating a slice of regexps and a slice of byte slices
	// containing the replacement strings.
	rp, r := compile(bl)

	// Find all the text files in the directory
	// This involves using the filepath.Glob function to find all files
	// in the current directory with the extension ".txt".
	tfs, err := findTextFiles()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Process each text file with the compiled blacklist
	// This involves reading each file into a byte slice, updating the
	// content of the byte slice using the regexps and replacement strings,
	// and then writing the updated byte slice back to the file.
	processFiles(tfs, rp, r)
}

// loadBlacklist takes no arguments and returns a map of strings to
// strings and an error value. The map is loaded from the "blacklist.json"
// file, and the error returned is non-nil if there's an error loading
// the blacklist.
//
// The function opens the "blacklist.json" file and decodes the file into
// the map. After decoding the map, the file is closed.
func loadBlacklist() (map[string]string, error) {
	f, err := os.Open("blacklist.json")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var m map[string]string
	err = json.NewDecoder(f).Decode(&m)
	return m, err
}

// compile takes a map of pattern strings to replacement strings
// and returns a slice of regexps and a slice of byte slices.
// The regexps are compiled from the patterns in the map.
// The byte slices are the replacement strings as byte slices.
func compile(bl map[string]string) ([]*regexp.Regexp, [][]byte) {
	var rs []*regexp.Regexp
	var rb [][]byte
	for p, r := range bl {
		// Compile the pattern string into a regex
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, nil
		}
		// Append the regex to the slice of regexps
		rs = append(rs, re)
		// Append the replacement string as a byte slice to the slice of byte slices
		rb = append(rb, []byte(r))
	}
	return rs, rb
}

// This function uses the filepath.Glob function to find all files in
// the current directory with the extension ".txt".  This function
// is used to find all text files in the current directory, which are
// then processed by the processFiles function.
func findTextFiles() ([]string, error) {
	f, err := filepath.Glob("*.txt")
	return f, err
}

// This function processes all the text files found by the findTextFiles
// function. It takes a slice of strings containing the names of the
// files to be processed, a slice of regexps compiled from the
// blacklist, and a slice of byte slices containing the replacement
// strings as byte slices.
//
// For each file, it reads the file into a byte slice, updates the
// content of the byte slice using the regexps and replacement strings,
// and then writes the updated byte slice back to the file.
func processFiles(names []string, re []*regexp.Regexp, rep [][]byte) {
	for _, name := range names {
		// Read the file into a byte slice
		data, err := os.ReadFile(name)
		if err != nil {
			// If there's an error, skip this file
			continue
		}

		// Update the content of the byte slice using the regexps and
		// replacement strings
		data = updateContent(data, re, rep)

		// Convert any years found in the text to the Buddhist Era
		data = convertYearsToBE(data)

		// Write the updated byte slice back to the file
		if err := os.WriteFile(name, data, os.ModePerm); err != nil {
			// If there's an error, skip this file
			continue
		}
	}
}

// This function takes a byte slice, a slice of regexps, and a slice of
// byte slices, and updates the byte slice using the regexps and
// replacement strings.
//
// For each regex in the slice of regexps, this function replaces all
// occurrences of the regex in the byte slice with the corresponding
// replacement string from the slice of byte slices.
//
// The updated byte slice is then returned.
func updateContent(data []byte, re []*regexp.Regexp, rep [][]byte) []byte {
	for i, r := range re {
		// Replace all occurrences of the regex in the byte slice with the
		// corresponding replacement string from the slice of byte slices
		data = r.ReplaceAll(data, rep[i])
	}
	return data
}

// This function takes a byte slice containing text, and converts any
// years found in the text to the Buddhist Era. This is done by using
// the regexp.MustCompile function to compile a regular expression matching
// any sequence of four digits, and then using the ReplaceAllFunc
// method to replace all matches with the corresponding year in the
// Buddhist Era.
//
// The replacement function takes a byte slice containing the matched
// text, converts it to an integer, and checks if the year is
// greater than or equal to 2400. If it is, it subtracts 543 from
// the year to convert it to the Buddhist Era, and returns the
// resulting year as a byte slice. If the year is less than 2400, it
// simply returns the original matched text.
func convertYearsToBE(c []byte) []byte {
	r := regexp.MustCompile(`(\d{4})`)
	return r.ReplaceAllFunc(c, func(m []byte) []byte {
		y, _ := strconv.Atoi(string(m))
		if y >= 2400 {
			return []byte(strconv.Itoa(y - 543))
		}
		return m
	})
}
