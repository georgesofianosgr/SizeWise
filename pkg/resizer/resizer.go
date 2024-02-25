package resizer

import (
	"fmt"
	"io"
	"os"

	"github.com/h2non/bimg"
)

func ResizeFile(inputPath, outputPath string, width, height int) error {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input image: %w", err)
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output image: %w", err)
	}

	err = Resize(inputFile, outputFile, width, height)
	if err != nil {
		return fmt.Errorf("failed to resize image: %w", err)
	}

	return nil
}

// Resize resizes an image to the specified width and height
func Resize(reader io.Reader, writer io.Writer, width, height int) error {
	buffer, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("failed to read input image: %w", err)
	}

	options := bimg.Options{
		Width:  width,
		Height: height,
	}

	newImage, err := bimg.Resize(buffer, options)
	if err != nil {
		return fmt.Errorf("failed to process image: %w", err)
	}

	_, err = writer.Write(newImage)
	if err != nil {
		return fmt.Errorf("failed to write image: %w", err)
	}
	return nil
}
