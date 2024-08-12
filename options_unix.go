//go:build linux || darwin

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

// WithSeparateProcessGroup set's the process group of the `xz` subprocess to its own, separate process group.  When
// the program using this library is started in a shell session, hitting CTRL+C will send an interrupt signal to both
// processes.  That means that the `xz` process terminates immediately without reading STDIN to its end, instead
// writes to the xzwriter fails with ERRPIPE.
//
// If the program using this library wants to handle the SIGINT gracefully, one needs to prevent the shell from sending
// the SIGINT to the xz subprocess also.  Running `xz` in a separate process group achieves that.
func WithSeparateProcessGroup() Option {
	return func(xz *XZWriter) error {
		xz.opts.separateProcessGroup = true

		return nil
	}
}

func WithNiceness(n int) Option {
	return func(xz *XZWriter) error {
		if n < 0 || n > 20 {
			return ErrOptionIllegal
		}

		xz.opts.niceness = n

		return nil
	}
}
