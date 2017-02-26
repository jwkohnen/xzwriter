/*
 * Copyright (c) 2017 Johannes Kohnen <wjkohnen@users.noreply.github.com>
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

package xz

import (
	"context"
	"io"
	"os/exec"

	ukxz "github.com/ulikunitz/xz"
)

// Writer is a WriteCloser that wraps a writer around an XZ compressor.
type Writer struct {
	cmd  *exec.Cmd
	pipe io.WriteCloser
}

// NewWriter returns an Writer, wrapping the writer w.
func NewWriter(ctx context.Context, w io.Writer) (*Writer, error) {
	xz := &Writer{}

	if isForceUKXZ(ctx) || xzPath == "" {
		uxz, err := ukxz.NewWriter(w)
		if err != nil {
			return nil, err
		}
		xz.pipe = uxz

		return xz, nil
	}

	cmd, pipe, err := xzWriteCmd(ctx, w)
	if err != nil {
		return nil, err
	}
	xz.cmd = cmd
	xz.pipe = pipe

	return xz, err
}

func xzWriteCmd(ctx context.Context, w io.Writer) (*exec.Cmd, io.WriteCloser, error) {
	cmd := exec.CommandContext(ctx, xzPath, "--quiet", "--compress",
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
func (xz *Writer) Write(p []byte) (n int, err error) {
	return xz.pipe.Write(p)
}

// Close implements the io.Closer interface.
func (xz *Writer) Close() error {
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
