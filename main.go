package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

// TextProc represents a text processing utility.
type TextProc struct{ R map[string]string }

// Process replaces words in the given text based on predefined replacements.
func (p TextProc) Process(t string) string {
	for w, r := range p.R {
		t = regexp.MustCompile(w).ReplaceAllString(t, r)
	}
	return t
}

// FileMgr manages file-related operations using a TextProc.
type FileMgr struct{ P TextProc }

// ProcessFile reads, processes, and writes the updated text to a file.
func (m FileMgr) ProcessFile(fname string) error {
	data, err := os.ReadFile(fname)
	if err != nil {
		return fmt.Errorf("error reading file %s: %v", fname, err)
	}

	ut := m.P.Process(string(data))
	if err := os.WriteFile(fname, []byte(ut), 0644); err != nil {
		return fmt.Errorf("error writing file %s: %v", fname, err)
	}

	fmt.Printf("Text in %s has been updated.\n", fname)
	return nil
}

// InitProc initializes a TextProc with replacements from a JSON file.
func InitProc() TextProc {
	file, err := os.Open("blacklist.json")
	if err != nil {
		panic(fmt.Errorf("error opening blacklist.json: %v", err))
	}
	defer file.Close()

	var r map[string]string
	if err := json.NewDecoder(file).Decode(&r); err != nil {
		panic(fmt.Errorf("error decoding blacklist.json: %v", err))
	}

	return TextProc{R: r}
}

func main() {
	txtFiles, err := filepath.Glob("*.txt")
	if err != nil {
		panic(err)
	}

	manager := FileMgr{P: InitProc()}
	for _, fname := range txtFiles {
		if err := manager.ProcessFile(fname); err != nil {
			fmt.Println(err)
		}
	}
}