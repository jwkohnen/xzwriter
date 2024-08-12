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
	"context"
	"errors"
	"io"
	"runtime"
	"testing"
)

func TestWithNiceness(t *testing.T) {
	xz, err := NewWithOptions(context.Background(), io.Discard, WithNiceness(20))
	switch runtime.GOOS {
	case "windows":
		if !errors.Is(err, ErrOptionIllegal) {
			t.Errorf("want %T, but got %T: %v", ErrOptionIllegal, err, err)
		}
	default:
		if err != nil {
			t.Fatalf("want no error, but got %v", err)
		}
	}

	if _, e := xz.Write([]byte("Hallo du da im Fernsehen!")); e != nil {
		t.Error(e)
	}

	if e := xz.Close(); e != nil {
		t.Error(e)
	}
}

// avoid annoying "no usage" warnings in IDE
var (
	_ = New
	_ = WithCompressLevel
	_ = WithExtreme
	_ = WithVerbose
	_ = WithSeparateProcessGroup
)

// asserts
var (
	_ io.WriteCloser = (*XZWriter)(nil)
)