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
	"io/ioutil"
	"os"
	"testing"
)

func doTestRead(ctx context.Context, t *testing.T) {
	f, err := os.Open("test.xz")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	r, err := NewReader(ctx, f)
	if err != nil {
		t.Fatal(err)
	}
	d, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(d, []byte("Hello, world!")) {
		t.Errorf("Want \"Hello, world!\", got %v", string(d))
	}
}

func doBenchmarkReader(ctx context.Context, b *testing.B) {
	for i := 0; i < b.N; i++ {
		func() {
			f, err := os.Open("test.xz")
			if err != nil {
				b.Fatal(err)
			}
			defer f.Close()
			r, err := NewReader(ctx, f)
			if err != nil {
				b.Fatal(err)
			}
			_, err = ioutil.ReadAll(r)
			if err != nil {
				b.Fatal(err)
			}
		}()
	}
}

func TestReader(t *testing.T) {
	doTestRead(context.Background(), t)
}

func TestReaderUKXZ(t *testing.T) {
	ctx := forceUKXZ(context.Background())
	doTestRead(ctx, t)
}

func BenchmarkReader(b *testing.B) {
	doBenchmarkReader(context.Background(), b)
}

func BenchmarkReaderUKXZ(b *testing.B) {
	ctx := forceUKXZ(context.Background())
	doBenchmarkReader(ctx, b)
}
