package xz

import (
	"context"
	"io"
	"os/exec"
)

type Reader struct {
	cmd  *exec.Cmd
	pipe io.ReadCloser
}

func NewReader(ctx context.Context, r io.Reader) (io.Reader, error) {
	if xzPath == "" {
		panic("no xz")
	}

	cmd := exec.CommandContext(ctx, xzPath, "-d", "-c", "-")
	cmd.Stdin = r
	pipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	return &Reader{cmd: cmd, pipe: pipe}, nil
}

func (x *Reader) Read(p []byte) (n int, err error) {
	n, err = x.pipe.Read(p)
	if err != nil {
		// Wait and cleanup process on any error including io.EOF
		errWait := x.cmd.Wait()
		if errWait != nil {
			err = errWait
		}
	}
	return n, err
}

var _ io.Reader = &Reader{}
