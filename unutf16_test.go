package unutf16_test

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nolotz/unutf16"
)

// TestUTF16LEToUTF8 tests conversion of UTF-16LE to UTF-8
func TestUTF16LEToUTF8(t *testing.T) {
	// UTF-16LE data (BOM + "hello")
	utf16leData := []byte{0xFF, 0xFE, 0x68, 0x00, 0x65, 0x00, 0x6C, 0x00, 0x6C, 0x00, 0x6F, 0x00}

	reader := bytes.NewReader(utf16leData)
	utf8Reader := unutf16.NewReader(reader)

	// Expected UTF-8 output
	expected := "hello"
	var output bytes.Buffer
	_, err := io.Copy(&output, utf8Reader)
	if err != nil {
		t.Fatalf("Error reading from UTF8 reader: %v", err)
	}

	if output.String() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, output.String())
	}
}

// TestUTF16BEToUTF8 tests conversion of UTF-16BE to UTF-8
func TestUTF16BEToUTF8(t *testing.T) {
	// UTF-16BE data (BOM + "hello")
	utf16beData := []byte{0xFE, 0xFF, 0x00, 0x68, 0x00, 0x65, 0x00, 0x6C, 0x00, 0x6C, 0x00, 0x6F}

	reader := bytes.NewReader(utf16beData)
	utf8Reader := unutf16.NewReader(reader)

	// Expected UTF-8 output
	expected := "hello"
	var output bytes.Buffer
	_, err := io.Copy(&output, utf8Reader)
	if err != nil {
		t.Fatalf("Error reading from UTF8 reader: %v", err)
	}

	if output.String() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, output.String())
	}
}

// TestNonUTF16Passthrough tests that non-UTF-16 data is passed through unmodified (e.g., UTF-8).
func TestNonUTF16Passthrough(t *testing.T) {
	// UTF-8 data (no BOM)
	utf8Data := []byte("hello world")

	reader := bytes.NewReader(utf8Data)
	utf8Reader := unutf16.NewReader(reader)

	// Expected UTF-8 output (should match input exactly)
	expected := "hello world"
	var output bytes.Buffer
	_, err := io.Copy(&output, utf8Reader)
	if err != nil {
		t.Fatalf("Error reading from UTF8 reader: %v", err)
	}

	if output.String() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, output.String())
	}
}

// TestPeekFailure simulates a failure during the Peek operation by using ErrorReader.
func TestPeekFailure(t *testing.T) {
	// Create an ErrorReader that triggers an error after 0 bytes (to simulate a peek failure)
	reader := new(errorReader)

	// Create the UTF-16 to UTF-8 reader
	utf8Reader := unutf16.NewReader(reader)

	// Attempt to read from the reader, which should fail due to the simulated peek error
	buffer := make([]byte, 10)
	_, err := utf8Reader.Read(buffer)
	if err == nil {
		t.Fatalf("Expected an error during Read, but got none")
	}

	assert.IsType(t, new(unutf16.BOMPeekError), err)
	assert.ErrorIs(t, err, simulatedError)
	assert.Equal(t, "failed to peek BOM: simulated read error", err.Error())
}

var simulatedError = errors.New("simulated read error")

type errorReader struct{}

func (e *errorReader) Read(p []byte) (int, error) {
	return 0, simulatedError
}
