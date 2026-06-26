package zipstream

import (
	"archive/zip"
	"fmt"
	"io"
)

// WriteFiltered writes a filtered zip archive to dst.
func WriteFiltered(
	dst io.Writer,
	files []*zip.File,
	keep Predicate,
	bufferSize int,
) error {
	if bufferSize <= 0 {
		bufferSize = defaultCopyBufferSize
	}

	writer := zip.NewWriter(dst)
	buffer := make([]byte, bufferSize)

	for _, file := range files {
		if !keep(file) {
			continue
		}
		header := RewriteRawHeader(&file.FileHeader)

		target, err := writer.CreateRaw(header)
		if err != nil {
			return fmt.Errorf("create raw entry %q: %w", file.Name, err)
		}

		source, err := file.OpenRaw()
		if err != nil {
			return fmt.Errorf("open raw entry %q: %w", file.Name, err)
		}

		if _, err := io.CopyBuffer(target, source, buffer); err != nil {
			return fmt.Errorf("copy raw entry %q: %w", file.Name, err)
		}
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("close zip writer: %w", err)
	}
	return nil
}

// RewriteRawHeader prepares a zip.FileHeader for CreateRaw.
func RewriteRawHeader(source *zip.FileHeader) *zip.FileHeader {
	clone := *source

	// We write known sizes directly into the local file header.
	clone.Flags &^= dataDescriptorFlag

	// Avoid carrying over extra records from the source archive.
	// This keeps output predictable and avoids conflicting metadata.
	clone.Extra = nil

	return &clone
}
