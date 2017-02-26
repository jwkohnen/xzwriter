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

type Reader struct {
	cmd    *exec.Cmd
	reader io.Reader
}

func NewReader(ctx context.Context, r io.Reader) (io.Reader, error) {
	xz := &Reader{}

	if isForceUKXZ(ctx) || xzPath == "" {
		uxz, err := ukxz.NewReader(r)
		if err != nil {
			return nil, err
		}
		xz.reader = uxz

		return xz, nil
	}

	cmd, reader, err := xzReadCmd(ctx, r)
	if err != nil {
		return nil, err
	}
	xz.cmd = cmd
	xz.reader = reader

	return xz, nil
}

func xzReadCmd(ctx context.Context, r io.Reader) (*exec.Cmd, io.Reader, error) {
	cmd := exec.CommandContext(ctx, xzPath, "--quiet", "--decompress",
		"--stdout", "-")
	cmd.Stdin = r

	pipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}

	err = cmd.Start()
	if err != nil {
		return nil, nil, err
	}

	return cmd, pipe, nil
}

func (x *Reader) Read(p []byte) (n int, err error) {
	n, err = x.reader.Read(p)
	if err != nil {
		if x.cmd != nil {
			// Wait and cleanup process on any error including io.EOF
			errWait := x.cmd.Wait()
			if errWait != nil {
				err = errWait
			}
		}
	}
	return n, err
}

var _ io.Reader = &Reader{}
