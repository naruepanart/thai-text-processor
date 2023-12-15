# Text Replacer

This simple Go program, `main.go`, is designed to process text files by replacing specified patterns with corresponding values. It can be particularly useful for automating text modifications across multiple files.

## How it Works

1. **Loading Replacements:**
   - The program starts by loading a set of replacements from a JSON file named `blacklist.json`. This file contains a map of patterns to replace and their corresponding replacements.

2. **Processing Text:**
   - For each text file in the current directory with a `.txt` extension, the program reads the content and applies the replacements defined in the loaded JSON file.
   - Additional replacements are defined within the code itself to showcase the flexibility of the tool.

3. **Saving Changes:**
   - The updated text is then written back to the original file, overwriting its content.

## How to Use

1. Ensure your text files are in the same directory as the `main.go` file.
2. Customize the `blacklist.json` file with your desired patterns and replacements.
3. Run the program using the command:

```go 
go build -ldflags="-s -w" main.go
```

## Additional Customization

Feel free to modify the code to add or customize replacements based on your specific needs. The program is currently set up to replace numeric ranges, clean up parentheses, and modify certain linguistic constructs.

## Disclaimer

This tool directly modifies the content of text files. Use it responsibly and make sure to have backups of your files before running the program, especially when experimenting with new replacement patterns.