/*
 * Copyright (c) 2021 Johannes Kohnen <jwkohnen-github@ko-sys.com>
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
	"context"
	"io"
	"os/exec"
	"strconv"
)

// XZWriter is a WriteCloser that wraps a writer around an XZ compressor.
type XZWriter struct {
	cmd  *exec.Cmd
	pipe io.WriteCloser
	opts options
}

// New returns an XZWriter, wrapping the writer w.
func New(w io.Writer) (*XZWriter, error) {
	return NewWithContext(context.Background(), w)
}

// NewWithContext returns an XZWriter, wrapping the writer w. The context may
// be used to cancel or timeout the external compressor process.
//
// The context can be used to kill the external process early. You still need to
// call Close() to clean up ressources. Alternatively you may call Close()
// prematurely.
func NewWithContext(ctx context.Context, w io.Writer) (*XZWriter, error) {
	return NewWithOptions(ctx, w)
}

// NewWithOptions returns an XZWriter, wrapping the writer w.  The context may
// be used to cancel or timeout the external compressor process.
//
// The context can be used to kill the external process early. You still need to
// call Close() to clean up ressources. Alternatively you may call Close()
// prematurely.
//
// The compressor process can be configured with options.
func NewWithOptions(ctx context.Context, w io.Writer, opts ...Option) (*XZWriter, error) {
	if ctx == nil {
		panic("nil Context")
	}

	xz := &XZWriter{
		// default options
		opts: options{
			compressLevel: Best,
		},
	}

	for _, opt := range opts {
		if err := opt(xz); err != nil {
			return nil, err
		}
	}

	xz.cmd = exec.CommandContext(ctx, "xz", xz.compileArgs()...)
	xz.cmd.Stdout = w

	if xz.opts.verboseWriter != nil {
		xz.cmd.Stderr = xz.opts.verboseWriter
	}

	if xz.opts.separateProcessGroup {
		xz.cmd.SysProcAttr = sysProcAttr()
	}

	var err error
	xz.pipe, err = xz.cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	err = xz.cmd.Start()
	if err != nil {
		return nil, err
	}

	return xz, err
}

// Write implements the io.Writer interface.
func (xz *XZWriter) Write(p []byte) (n int, err error) {
	return xz.pipe.Write(p)
}

// Close implements the io.Closer interface.
func (xz *XZWriter) Close() error {
	errPipe := xz.pipe.Close()
	errWait := xz.cmd.Wait()
	if errPipe != nil {
		return errPipe
	}
	return errWait
}

func (xz *XZWriter) compileArgs() []string {
	compressLevel := "-" + strconv.Itoa(xz.opts.compressLevel)

	args := []string{"--compress", "--stdout", compressLevel}

	if xz.opts.extreme {
		args = append(args, "--extreme")
	}

	if xz.opts.verboseWriter != nil {
		args = append(args, "--verbose")
	} else {
		args = append(args, "--quiet")
	}

	return append(args, "--", "-")
}

var _ io.WriteCloser = &XZWriter{} // assert
