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
	"bytes"
	"context"
	"io"
	"os/exec"

	ukxz "github.com/ulikunitz/xz"
)

var (
	xzPath = findXZ()

	// type asserts
	_ io.WriteCloser = &Writer{}
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
	if !bytes.Contains(out, []byte("<http://tukaani.org/xz/>")) {
		return ""
	}
	return path
}

type modeKey struct {
	name string
}

var modeForceUKXZ = &modeKey{"force UKXZ"}

func forceUKXZ(ctx context.Context) context.Context {
	return context.WithValue(ctx, modeForceUKXZ, true)
}

func isForceUKXZ(ctx context.Context) bool {
	v, _ := ctx.Value(modeForceUKXZ).(bool)
	return v
}
