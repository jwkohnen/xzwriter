package xzwriter

import (
	"errors"
	"io"
)

type Option func(*XZWriter) error

var (
	ErrOptionIllegal = errors.New("option illegal")
)

const (
	Fast    = 0
	Default = 6
	Best    = 9
)

// WithCompressLevel sets the compression level between 0 and 9.  The constants `Fast`, `Default` and `Best` correspond
// to the flags `--fast`, `--default` and `--best`.
func WithCompressLevel(l int) Option {
	return func(xz *XZWriter) error {
		if l < Fast || l > Best {
			return ErrOptionIllegal
		}

		xz.opts.compressLevel = l

		return nil
	}
}

// WithExtreme set the `--extreme` flag.
func WithExtreme() Option {
	return func(xz *XZWriter) error {
		xz.opts.extreme = true

		return nil
	}
}

// WithVerbose sets verbosity and takes a writer that will be connected to STDERR of the xz subprocess.  This provides
// a nice progress output to look at.
func WithVerbose(stderr io.Writer) Option {
	return func(xz *XZWriter) error {
		xz.opts.verboseWriter = stderr

		return nil
	}
}

type options struct {
	compressLevel        int
	extreme              bool
	verboseWriter        io.Writer
	separateProcessGroup bool
}
