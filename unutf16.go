package unutf16

import (
	"bufio"
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
	// Create a buffered reader to ensure we can "peek" into the stream for BOM detection
	var bufReader *bufio.Reader
	if br, ok := r.source.(*bufio.Reader); ok {
		bufReader = br
	} else {
		bufReader = bufio.NewReader(r.source)
	}

	// Read the first 2 bytes to check for BOM
	bom, err := bufReader.Peek(2) // Peek allows us to see the first 2 bytes without consuming them
	if err != nil && err != io.EOF {
		return fmt.Errorf("failed to read BOM: %v", err)
	}

	// Detect BOM and create the appropriate decoder
	var decoder io.Reader
	if len(bom) >= 2 && bom[0] == 0xFF && bom[1] == 0xFE {
		// UTF-16 Little Endian
		decoder = transform.NewReader(bufReader, unicode.UTF16(unicode.LittleEndian, unicode.UseBOM).NewDecoder())
	} else if len(bom) >= 2 && bom[0] == 0xFE && bom[1] == 0xFF {
		// UTF-16 Big Endian
		decoder = transform.NewReader(bufReader, unicode.UTF16(unicode.BigEndian, unicode.UseBOM).NewDecoder())
	} else {
		decoder = bufReader
	}

	// Assign the decoder to the reader
	r.decoder = decoder
	return nil
}
