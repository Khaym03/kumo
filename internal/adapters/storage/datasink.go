package sinks

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// HTMLFileSink implements the ports.DataSink interface to save data to HTML files.
type HTMLFileSink struct {
	dir string
}

// NewHTMLFileSink creates a new instance of HTMLFileSink.
func NewHTMLFileSink(dir string) (*HTMLFileSink, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create output directory: %w", err)
	}
	return &HTMLFileSink{dir: dir}, nil
}

// Store saves the content to a file with the given ID and a .html extension.
func (s *HTMLFileSink) Store(id string, data io.Reader) error {
	filePath := filepath.Join(s.dir, fmt.Sprintf("%s.html", id))
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create HTML file %s: %w", filePath, err)
	}
	defer file.Close()

	if _, err := io.Copy(file, data); err != nil {
		return fmt.Errorf("failed to write to HTML file: %w", err)
	}

	return nil
}
