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
	f, err := os.Open("blacklist.json")
	if err != nil {
		panic(fmt.Errorf("open blacklist.json: %v", err))
	}
	defer f.Close()

	bl := make(map[string]string)
	if err := json.NewDecoder(f).Decode(&bl); err != nil {
		panic(fmt.Errorf("decode blacklist.json: %v", err))
	}

	rx := make([]*regexp.Regexp, 0, len(bl))
	rp := make([][]byte, 0, len(bl))
	for p, r := range bl {
		re, err := regexp.Compile(p)
		if err != nil {
			panic(fmt.Errorf("compile regex %s: %v", p, err))
		}
		rx = append(rx, re)
		rp = append(rp, []byte(r))
	}

	files, err := filepath.Glob("*.txt")
	if err != nil {
		panic(fmt.Errorf("glob: %v", err))
	}

	for _, fn := range files {
		file, err := os.OpenFile(fn, os.O_RDWR, os.ModePerm)
		if err != nil {
			fmt.Println(fmt.Errorf("open file %s: %v", fn, err))
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
				fmt.Println(fmt.Errorf("read file %s: %v", fn, err))
				break
			}

			txt := buf[:n]
			for i, re := range rx {
				txt = re.ReplaceAll(txt, rp[i])
			}

			updated = append(updated, txt...)
		}

		if err := os.Truncate(fn, 0); err != nil {
			fmt.Println(fmt.Errorf("truncate file %s: %v", fn, err))
			continue
		}
		if _, err := file.WriteAt(updated, 0); err != nil {
			fmt.Println(fmt.Errorf("write file %s: %v", fn, err))
			continue
		}

		fmt.Printf("Text in %s updated.\n", fn)
	}
}