package xz

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"testing"
)

const msg = "Hello, world!"

func doTestWrite(ctx context.Context, t *testing.T) {
	pr, pw := io.Pipe()

	xw, err := NewWriter(ctx, pw)
	if err != nil {
		t.Fatal(err)
	}
	xr, err := NewReader(ctx, pr)
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		n, err := fmt.Fprint(xw, msg)
		if err != nil {
			t.Fatal(err)
		}
		if n != len(msg) {
			t.Errorf("Expected to write %d bytes, wrote %d", len(msg), n)
		}
		err = xw.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()
	d, err := ioutil.ReadAll(xr)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(d, []byte(msg)) {
		t.Errorf("Want \"%s\", got %s", msg, string(d))
	}
}

func TestWriter(t *testing.T) {
	doTestWrite(context.Background(), t)
}

func TestWriterUKXZ(t *testing.T) {
	doTestWrite(forceUKXZ(context.Background()), t)
}
