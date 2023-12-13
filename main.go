package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

// TextProcessor interface defines the method for processing text.
type TextProcessor interface {
	ProcessText(text string) string
}

// ThaiTextProcessor implements the TextProcessor interface for Thai text processing.
type ThaiTextProcessor struct {
	WordReplacements map[string]string
}

// ProcessText replaces words in the text based on the provided replacements.
func (p ThaiTextProcessor) ProcessText(text string) string {
	for word, replacement := range p.WordReplacements {
		text = regexp.MustCompile(word).ReplaceAllString(text, replacement)
	}
	return text
}

// FileReader interface defines the method for reading a file.
type FileReader interface {
	ReadFile(filename string) ([]byte, error)
}

// OSFileReader implements the FileReader interface using the OS package.
type OSFileReader struct{}

// ReadFile reads the content of the file using the OS package.
func (r OSFileReader) ReadFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

// FileWriter interface defines the method for writing to a file.
type FileWriter interface {
	WriteFile(filename string, data []byte) error
}

// OSFileWriter implements the FileWriter interface using the OS package.
type OSFileWriter struct{}

// WriteFile writes the data to the file using the OS package.
func (w OSFileWriter) WriteFile(filename string, data []byte) error {
	return os.WriteFile(filename, data, 0644)
}

// FileManager manages the processing and updating of text files.
type FileManager struct {
	TextProcessor TextProcessor
	FileReader    FileReader
	FileWriter    FileWriter
}

// ProcessFile reads, processes, and writes a file using FileManager.
func (m FileManager) ProcessFile(filename string) error {
	data, err := m.FileReader.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error reading file %s: %v", filename, err)
	}

	updatedText := m.TextProcessor.ProcessText(string(data))

	err = m.FileWriter.WriteFile(filename, []byte(updatedText))
	if err != nil {
		return fmt.Errorf("error writing file %s: %v", filename, err)
	}

	fmt.Printf("Text in %s has been updated.\n", filename)

	return nil
}

// getCurrentDirectory returns the current working directory.
func getCurrentDirectory() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return dir
}

// initializeThaiBlacklist initializes a map of Thai words and their replacements.
func initializeThaiBlacklist() map[string]string {
	return map[string]string{
		" ๆ":                     "ๆ",
		"เรียกว่า":               " เรียกว่า ",
		"และ":                    " และ",
		"ซึ่ง":                   " ซึ่ง",
		"สามารถ":                 " สามารถ",
		"ฟังก์ชัน":               "ฟังก์ชั่น",
		"ลูชัน":                  "ลูชั่น",
		"ต่อๆ ไป":                "ต่อๆไป",
		"อัลกอริทึม":             "อัลกอริทึ่ม",
		"อัลกอริธึม":             "อัลกอริธึ่ม",
		"กาแลคซี":                "กาแล็คซี่",
		"กาแล็กซี":               "กาแล็กซี่",
		"แอนโดรเมดา":             "แอนโดรเมด้า",
		"สสาร":                   "สะสาร",
		"ไดนามิก":                "ไดนามิค",
		"เอนโทรปี":               "เอนโทรปี้",
		"ยูโรปา":                 "ยูโรป้า",
		"ฮับเบิล":                "ฮับเบิ้ล",
		"ควอนตัม":                "ควอนตั้ม",
		"ละลาย":                  "ละลาย ",
		"ความจำเป็น":             "ความ จำเป็น",
		"อัปโหลด":                "อัพโหลด",
		"พลาสมา":                 "พลาสม่า",
		"มีนัย":                  "มีนัยะ",
		"ยักษ์":                  "ยัก",
		"แอริโซนา":               "แอริโซน่า",
		"ลิเธียม":                "ลิเธี่ยม",
		" กม. ":                  " กิโลเมตร ",
		"อาร์เทมิส":              "อาร์ทิมิส",
		" CO2 ":                  " C O 2 ",
		" EVs ":                  " E V ",
		"GPS":                    "G P S",
		"Square Kilometer Array": "สแควร์ กิโลมิเตอร์ อาร์เรย์",
		"SpaceX":                 "สเปซเอ็กซ์",
		"Blue Origin":            "บลูออริจิ้น",
		"Curiosity":              "คิวริออซิตี้",
		"Galileo":                "กาลิเลโอ",
		"STEM":                   "สะเต็ม",
		"NASA":                   "นาซ่า",
		"Perseverance":           "เพอซะเวียแร้น",
		"°F":                     "องศาฟาเรนไฮต์",
		"°C":                     "องศาเซลเซียส",
		":":                      "",
		"โดยรวมแล้ว ":            "",
		"โดยสรุป":                "",
	}
}

func main() {
	scriptDir := getCurrentDirectory()

	txtFiles, err := filepath.Glob(filepath.Join(scriptDir, "*.txt"))
	if err != nil {
		panic(err)
	}

	blacklist := initializeThaiBlacklist()

	fileManager := FileManager{
		TextProcessor: ThaiTextProcessor{WordReplacements: blacklist},
		FileReader:    OSFileReader{},
		FileWriter:    OSFileWriter{},
	}

	for _, fileName := range txtFiles {
		if err := fileManager.ProcessFile(fileName); err != nil {
			fmt.Println(err)
		}
	}
}
