package resizer

import (
	"os"
	"path/filepath"
	"testing"
)

// CreateTempFile creates a temporary file and returns its path
func CreateTempFile(pattern string) (string, error) {
	file, err := os.CreateTemp("/tmp", pattern)
	return file.Name(), err
}

func TestResizeFile(t *testing.T) {
	inputPath := filepath.Join("testdata", "test.jpeg")
	file, err := os.CreateTemp("/tmp", "*.jpeg")
	if err != nil {
		t.Fatalf("got error %q", err)
	}
	outputPath := file.Name()
	file.Close()
	defer os.Remove(outputPath)

	err = ResizeFile(inputPath, outputPath, 100, 100)
	if err != nil {
		t.Fatalf("got error %q", err)
	}
}
