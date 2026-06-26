package zipstream

import (
	"archive/zip"
	"io"
)

const (
	defaultCopyBufferSize = 256 * 1024
	dataDescriptorFlag    = 0x0008
)

// ReaderAtCloser is an interface that combines io.ReaderAt and io.Closer.
//
// It is required because zip.NewReader requires an io.ReaderAt, but we
// also need to close the source when done.
type ReaderAtCloser interface {
	io.ReaderAt
	io.Closer
}

// Predicate is a function that returns true if a zip.File should be
// kept in the output archive.
type Predicate func(file *zip.File) bool

// StreamFiltered streams a new zip archive containing only entries
// for which keep(file) returns true.
//
// The source is owned by this function and will be closed when
// streaming ends, so the caller should not close it.
func StreamFiltered(
	source ReaderAtCloser,
	size int64,
	keep Predicate,
	optionList ...Option,
) (io.ReadCloser, int64, error) {
	opts := options{
		bufferSize: defaultCopyBufferSize,
	}
	for _, option := range optionList {
		option(&opts)
	}

	reader, err := zip.NewReader(source, size)
	if err != nil {
		source.Close()
		return nil, -1, err
	}
	var resultSize int64 = -1

	if archiveSize, ok := ArchiveSizeFiltered(reader.File, keep); ok {
		resultSize = archiveSize
	}

	pipeReader, pipeWriter := io.Pipe()

	go func() {
		defer source.Close()
		err := WriteFiltered(pipeWriter, reader.File, keep, opts.bufferSize)
		pipeWriter.CloseWithError(err)
	}()

	return pipeReader, resultSize, nil
}
