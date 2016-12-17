/*
 * Copyright (c) 2016 Wolfgang Johannes Kohnen <wjkohnen@users.noreply.github.com>
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
// XZ compressor.
//
// Uses the Tukaani XZ tool in $PATH if available. See the XZ Utils home page:
// <http://tukaani.org/xz/>. If the XZ Utils are not found, XZWriter will
// fall back to the Go native implementation by Ulrich Kunitz
// <https://github.com/ulikunitz/xz>.
//
// WARNING: The xz writer by Ulrich Kunitz is alpha quality software and may
// destroy your data. This packages hides the fact which compressor will be
// used. If you have high expections on integrity, prepend a hasher to the
// writer chain, re-read written data and compare!
package xzwriter

import (
	"io"
	"os/exec"
	"strings"

	ukxz "github.com/ulikunitz/xz"
)

// XZWriter is a WriteCloser that wraps a writer around an XZ compressor.
type XZWriter struct {
	cmd  *exec.Cmd
	pipe io.WriteCloser
}

// New returns an XZWriter, wrapping the writer w.
func New(w io.Writer) (*XZWriter, error) {
	xz := &XZWriter{}
	var err error

	if xzPath == "" {
		uxz, err := ukxz.NewWriter(w)
		if err != nil {
			return nil, err
		}
		xz.pipe = uxz

		return xz, nil
	}

	cmd, pipe, err := xzCmd(w)
	if err != nil {
		return nil, err
	}
	xz.cmd = cmd
	xz.pipe = pipe

	return xz, err
}

func xzCmd(w io.Writer) (*exec.Cmd, io.WriteCloser, error) {
	cmd := exec.Command(xzPath, "--quiet", "--compress",
		"--stdout", "--best", "-")
	cmd.Stdout = w

	pipe, err := cmd.StdinPipe()
	if err != nil {
		return nil, nil, err
	}

	err = cmd.Start()
	if err != nil {
		return nil, nil, err
	}

	return cmd, pipe, nil
}

// Write implements the io.Writer interface.
func (xz *XZWriter) Write(p []byte) (n int, err error) {
	return xz.pipe.Write(p)
}

// Close implements the io.Closer interface.
func (xz *XZWriter) Close() error {
	var errPipe, errWait error

	errPipe = xz.pipe.Close()
	if xz.cmd != nil {
		errWait = xz.cmd.Wait()
	}
	if errPipe != nil {
		return errPipe
	}
	return errWait
}

var (
	xzPath = findXZ()

	// type asserts
	_ io.WriteCloser = &XZWriter{}
	_ io.WriteCloser = &ukxz.Writer{}
)

func findXZ() string {
	path, err := exec.LookPath("xz")
	if err != nil {
		return ""
	}
	cmd := exec.Command(path, "--help")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}
	if !strings.Contains(string(out), "<http://tukaani.org/xz/>") {
		return ""
	}
	return path
}
