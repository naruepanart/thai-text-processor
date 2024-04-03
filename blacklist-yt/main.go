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
	b, err := os.Open("blacklist.json")
	if err != nil {
		panic(fmt.Errorf("open blacklist.json: %v", err))
	}
	defer b.Close()

	r := make(map[string]string)
	if err := json.NewDecoder(b).Decode(&r); err != nil {
		panic(fmt.Errorf("decode blacklist.json: %v", err))
	}

	p := make([]*regexp.Regexp, 0, len(r))
	v := make([]string, 0, len(r))
	for pattern, value := range r {
		p = append(p, regexp.MustCompile(pattern))
		v = append(v, value)
	}

	t, err := filepath.Glob("*.txt")
	if err != nil {
		panic(err)
	}

	for _, f := range t {
		file, err := os.Open(f)
		if err != nil {
			fmt.Println(fmt.Errorf("open file %s: %v", f, err))
			continue
		}
		defer file.Close()

		var updated []byte
		buf := make([]byte, 1024)

		for {
			n, err := file.Read(buf)
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Println(fmt.Errorf("read file %s: %v", f, err))
				break
			}

			text := buf[:n]
			for i, pattern := range p {
				text = pattern.ReplaceAll(text, []byte(v[i]))
			}

			updated = append(updated, text...)
		}

		if err = os.WriteFile(f, updated, os.ModePerm); err != nil {
			fmt.Println(fmt.Errorf("write file %s: %v", f, err))
			continue
		}

		fmt.Printf("Text in %s updated.\n", f)
	}
}