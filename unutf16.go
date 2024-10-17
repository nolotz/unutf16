package unutf16

import (
	"bytes"
	"fmt"
	"io"

	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

// NewReader initializes a new Reader that wraps an existing io.Reader.
// This function prepares the Reader for converting UTF-16 encoded data to UTF-8,
// but does not start decoding until the first Read call is made.
// Returns a new Reader that wraps the provided io.Reader and handles UTF-16 to UTF-8 conversion.
func NewReader(r io.Reader) *Reader {
	return &Reader{
		source:  r,
		decoder: nil,
	}
}

// Reader is a custom io.Reader that wraps an existing io.Reader (source)
// and optionally converts UTF-16 encoded data into UTF-8.
// The decoder field is an internal io.Reader that handles the UTF-16 to UTF-8 conversion.
// If the source is already UTF-8 or doesn't require conversion, the decoder equals source.
type Reader struct {
	source  io.Reader // Underlying source reader (UTF-16 encoded)
	decoder io.Reader // Decoder that will handle the conversion from UTF-16 to UTF-8
}

// Read implements the io.Reader interface.
// It lazily initializes the decoder on the first read, then streams the converted content.
func (r *Reader) Read(p []byte) (int, error) {
	// Lazy initialization: perform BOM detection and setup the decoder on the first read call
	if r.decoder == nil {
		err := r.initialize()
		if err != nil {
			return 0, err
		}
	}

	// Now delegate the Read call to the decoder, which handles UTF-16 to UTF-8 conversion
	return r.decoder.Read(p)
}

// initialize sets up the decoder by detecting the BOM and initializing the appropriate transform.Reader.
func (r *Reader) initialize() error {
	bom := make([]byte, 2)
	// Read the first 2 bytes to check for BOM
	_, err := r.source.Read(bom)
	if err != nil && err != io.EOF {
		return &BOMPeekError{
			Cause: err,
		}
	}

	// Stitch everything back again
	newReader := io.MultiReader(bytes.NewReader(bom), r.source)

	// Detect BOM and create the appropriate decoder
	var decoder io.Reader
	if len(bom) >= 2 && bom[0] == 0xFF && bom[1] == 0xFE {
		// UTF-16 Little Endian
		decoder = transform.NewReader(newReader, unicode.UTF16(unicode.LittleEndian, unicode.UseBOM).NewDecoder())
	} else if len(bom) >= 2 && bom[0] == 0xFE && bom[1] == 0xFF {
		// UTF-16 Big Endian
		decoder = transform.NewReader(newReader, unicode.UTF16(unicode.BigEndian, unicode.UseBOM).NewDecoder())
	} else {
		decoder = newReader
	}

	// Assign the decoder to the reader
	r.decoder = decoder
	return nil
}

// BOMPeekError is a custom error type that represents an error encountered
// while attempting to peek the Byte Order Mark (BOM) from an input stream.
// This error wraps the original error (`Cause`) that occurred during the peek operation.
type BOMPeekError struct {
	Cause error
}

// Error implements the error interface for BOMPeekError.
// Returns a formatted error message that includes the underlying cause of the error.
//
// Example error message:
//
//	"failed to peek BOM: unexpected EOF"
func (e *BOMPeekError) Error() string {
	return fmt.Sprintf("failed to peek BOM: %v", e.Cause)
}

// Unwrap allows the BOMPeekError to expose the underlying error that caused the failure.
// This can be used to retrieve the original error when handling multiple layers of errors.
//
// Example usage:
//
//	if errors.Is(err, io.EOF) { ... }  // Allows matching against the wrapped error.
func (e *BOMPeekError) Unwrap() error {
	return e.Cause
}
