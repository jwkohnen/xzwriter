/*
 * Copyright (c) 2024 Johannes Kohnen <jwkohnen-github@ko-sys.com>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package xzwriter provides a WriteCloser XZWriter that pipes through an
// external XZ compressor.
//
// Expects the Tukaani XZ tool in $PATH. See the XZ Utils home page:
// <http://tukaani.org/xz/>

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
// nice progress output to look at.
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
	niceness             int
}
