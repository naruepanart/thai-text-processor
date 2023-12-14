package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

type TextProc struct{ Replacements map[string]string }

func (p TextProc) Process(text string) string {
	for word, replacement := range p.Replacements {
		text = regexp.MustCompile(word).ReplaceAllString(text, replacement)
	}
	return text
}

type FileMgr struct{ Processor TextProc }

func (m FileMgr) ProcessFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error reading %s: %v", filename, err)
	}

	updatedText := m.Processor.Process(string(data))
	if err := os.WriteFile(filename, []byte(updatedText), 0644); err != nil {
		return fmt.Errorf("error writing %s: %v", filename, err)
	}

	fmt.Printf("Text in %s updated.\n", filename)
	return nil
}

func InitProc() TextProc {
	file, err := os.Open("blacklist.json")
	if err != nil {
		panic(fmt.Errorf("error opening blacklist.json: %v", err))
	}
	defer file.Close()

	var replacements map[string]string
	if err := json.NewDecoder(file).Decode(&replacements); err != nil {
		panic(fmt.Errorf("error decoding blacklist.json: %v", err))
	}

	replacements[`(\d+)-(\d+)`] = `${1}ถึง${2}`
	replacements[`(\p{L}+)\s+\(`] = `${1} หรือ `
	replacements[`(\s*\))`] = ``

	return TextProc{Replacements: replacements}
}

func main() {
	textFiles, err := filepath.Glob("*.txt")
	if err != nil {
		panic(err)
	}

	manager := FileMgr{Processor: InitProc()}
	for _, filename := range textFiles {
		if err := manager.ProcessFile(filename); err != nil {
			fmt.Println(err)
		}
	}
}