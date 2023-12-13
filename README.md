# Thai Text Processor

[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

## Overview

Thai Text Processor is a command-line tool for processing Thai text files. It utilizes a customizable word replacement mechanism to update and modify the content of text files.

## Features

- Thai language text processing
- Customizable word replacement using a blacklist
- Command-line interface for easy integration into scripts or workflows

## Installation

Clone the repository to your local machine and navigate to the project directory.

```bash
git clone https://github.com/naruepanart/thai-text-processor.git
cd thai-text-processor
```

## Usage

1. Ensure that you have Thai text files (with a `.txt` extension) in the same directory as the script.
2. Run the script using the following command:

```bash
go run main.go
```

The script will process each Thai text file in the directory, applying word replacements based on the specified blacklist.

## Customization

### Adding or Modifying Word Replacements

To customize word replacements, open the `initializeThaiBlacklist` function in the `main.go` file. Add or modify entries as needed. The format is `"original word": "replacement"`.

Example:

```go
"example": "custom-replacement",
```

## Contributing

If you would like to contribute to the project, please follow these steps:

1. Fork the repository.
2. Create a new branch for your feature or bug fix.
3. Make your changes and commit them with descriptive commit messages.
4. Push your changes to your fork.
5. Open a pull request, explaining the purpose and benefits of your changes.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

Special thanks to the contributors and libraries that make this project possible.

- [Regular Expressions Package](https://golang.org/pkg/regexp/)

## Contact

For any questions or concerns, feel free to reach out:

- Email: naruepanart1201@gmail.com