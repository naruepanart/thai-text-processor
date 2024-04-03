package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Node struct {
	children map[byte]*Node
	isEnd    bool
	value    string
}

func main() {
	b, e := os.Open("blacklist.json")
	if e != nil {
		panic(fmt.Errorf("failed to open blacklist.json: %v", e))
	}
	defer b.Close()

	var r map[string]string
	if e := json.NewDecoder(b).Decode(&r); e != nil {
		panic(fmt.Errorf("failed to decode blacklist.json: %v", e))
	}

	root := &Node{children: make(map[byte]*Node)}
	for p, v := range r {
		addToTrie(root, p, v)
	}

	t, e := filepath.Glob("*.txt")
	if e != nil {
		panic(e)
	}

	for _, f := range t {
		file, e := os.Open(f)
		if e != nil {
			fmt.Println(fmt.Errorf("failed to open file %s: %v", f, e))
			continue
		}
		defer file.Close()

		var u []byte
		buf := make([]byte, 1024)

		for {
			n, e := file.Read(buf)
			if e == io.EOF {
				break
			}
			if e != nil {
				fmt.Println(fmt.Errorf("failed to read file %s: %v", f, e))
				break
			}

			text := buf[:n]
			updatedText := process(root, text)
			u = append(u, updatedText...)
		}

		if e := os.WriteFile(f, u, os.ModePerm); e != nil {
			fmt.Println(fmt.Errorf("failed to write file %s: %v", f, e))
			continue
		}
		fmt.Printf("Text in %s updated.\n", f)
	}
}

func addToTrie(root *Node, pattern, value string) {
	node := root
	for i := 0; i < len(pattern); i++ {
		char := pattern[i]
		if node.children == nil {
			node.children = make(map[byte]*Node)
		}
		if _, ok := node.children[char]; !ok {
			node.children[char] = &Node{}
		}
		node = node.children[char]
	}
	node.isEnd = true
	node.value = value
}

func process(root *Node, text []byte) []byte {
	var updated []byte
	node := root
	start := 0
	for i := 0; i < len(text); i++ {
		char := text[i]
		if node.children[char] != nil {
			node = node.children[char]
			if node.isEnd {
				updated = append(updated, []byte(node.value)...)
				start = i + 1
				node = root
			}
		} else {
			updated = append(updated, text[start:i+1]...)
			start = i + 1
			node = root
		}
	}
	if start < len(text) {
		updated = append(updated, text[start:]...)
	}
	return updated
}