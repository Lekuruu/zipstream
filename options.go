package zipstream

type options struct {
	bufferSize int
}

// Option is a function that modifies the options for the zipstream.
type Option func(*options)

// WithBufferSize returns an Option that sets the buffer size for the zipstream.
func WithBufferSize(size int) Option {
	return func(opts *options) {
		opts.bufferSize = size
	}
}
