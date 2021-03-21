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

func WithCompressLevel(l int) Option {
	return func(xz *XZWriter) error {
		if l < 0 || l > 9 {
			return ErrOptionIllegal
		}

		xz.opts.compressLevel = l

		return nil
	}
}

func WithExtreme() Option {
	return func(xz *XZWriter) error {
		xz.opts.extreme = true

		return nil
	}
}

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
