package xz

import (
	"bytes"
	"context"
	"io/ioutil"
	"os"
	"testing"
)

func TestRead(t *testing.T) {
	f, err := os.Open("test.xz")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	r, err := NewReader(context.Background(), f)
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

func BenchmarkReader(b *testing.B) {
	for i := 0; i < b.N; i++ {
		func() {
			f, err := os.Open("test.xz")
			if err != nil {
				b.Fatal(err)
			}
			defer f.Close()
			r, err := NewReader(context.Background(), f)
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
