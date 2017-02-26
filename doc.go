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
package xz
