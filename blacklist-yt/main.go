package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
)

func main() {
	err := processFiles(); 
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

func processFiles() error {
	bl, err := loadBlacklist("blacklist.json")
	if err != nil {
		return fmt.Errorf("loading blacklist: %v", err)
	}

	rx, rp, err := compileRegex(bl)
	if err != nil {
		return fmt.Errorf("compiling regex: %v", err)
	}

	files, err := filepath.Glob("*.txt")
	if err != nil {
		return fmt.Errorf("glob: %v", err)
	}

	for _, fn := range files {
		if err := processFile(fn, rx, rp); err != nil {
			fmt.Printf("Error processing %s: %v\n", fn, err)
		}
	}

	return nil
}

func loadBlacklist(filename string) (map[string]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var bl map[string]string
	if err := json.NewDecoder(f).Decode(&bl); err != nil {
		return nil, err
	}
	return bl, nil
}

func compileRegex(bl map[string]string) ([]*regexp.Regexp, [][]byte, error) {
	var rx []*regexp.Regexp
	var rp [][]byte

	for p, r := range bl {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, nil, err
		}
		rx = append(rx, re)
		rp = append(rp, []byte(r))
	}

	return rx, rp, nil
}

func processFile(filename string, rx []*regexp.Regexp, rp [][]byte) error {
	file, err := os.OpenFile(filename, os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	var updated []byte
	buf := make([]byte, 1024)

	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		txt := buf[:n]
		for i, re := range rx {
			txt = re.ReplaceAll(txt, rp[i])
		}

		updated = append(updated, txt...)
	}

	if err := os.Truncate(filename, 0); err != nil {
		return err
	}
	if _, err := file.WriteAt(updated, 0); err != nil {
		return err
	}

	fmt.Printf("Text in %s updated.\n", filename)
	return nil
}