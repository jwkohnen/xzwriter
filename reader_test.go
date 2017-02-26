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
	ctx := ForceUKXZ(context.Background())
	doTestRead(ctx, t)
}

func BenchmarkReader(b *testing.B) {
	doBenchmarkReader(context.Background(), b)
}

func BenchmarkReaderUKXZ(b *testing.B) {
	ctx := ForceUKXZ(context.Background())
	doBenchmarkReader(ctx, b)
}
