# Lazy UTF-16 to UTF-8 Reader
[![GoDoc](https://img.shields.io/badge/godoc-reference-5272B4)](https://pkg.go.dev/mod/github.com/nolotz/unutf16)
[![Go Report Card](https://goreportcard.com/badge/github.com/nolotz/unutf16)](https://goreportcard.com/report/github.com/nolotz/unutf16)

This Go library provides a custom io.Reader that lazily converts UTF-16 encoded text streams into UTF-8. The reader supports automatic detection of Byte Order Mark (BOM) and converts both UTF-16 Little Endian (LE) and Big Endian (BE) to UTF-8, making it ideal for handling mixed or unknown UTF-16 encoded data sources.

It only initializes the decoding process when the first Read is called, ensuring efficient and delayed BOM detection.

### Features
- **Lazy Initialization:** BOM detection and decoder setup only happen upon the first read.
- **Supports UTF-16LE and UTF-16BE:** Automatically detects the endianness based on the BOM.
- **Streaming Support:** Works with io.Reader, making it memory-efficient for large files or streams.
- **Seamless Integration:** Can be used just like any other io.Reader in Go.

## Installation

```bash
go get github.com/nolotz/unutf16
```

## Usage

Here’s an example of how to use the UTF-16 to UTF-8 reader in your Go code:

```go
package main

import (
    "bytes"
    "fmt"
    "io"

    "github.com/nolotz/unutf16"
)

func main() {
    // Example UTF-16LE encoded data (with BOM)
    utf16leData := []byte{0xFF, 0xFE, 0x68, 0x00, 0x65, 0x00, 0x6C, 0x00, 0x6C, 0x00, 0x6F, 0x00} // "hello"

    reader := bytes.NewReader(utf16leData)

    // Create a new UTF-16 to UTF-8 reader
    utf8Reader := unutf16.NewReader(reader)

    // Read the converted UTF-8 data
    utf8Output := new(bytes.Buffer)
    _, err := io.Copy(utf8Output, utf8Reader)
    if err != nil {
        panic(err)
    }

    fmt.Println(utf8Output.String()) // Output: "hello"
}
```