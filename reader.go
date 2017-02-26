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
	} else {
		var err error
		xz.cmd, xz.reader, err = xzReadCmd(ctx, r)
		if err != nil {
			return nil, err
		}
	}

	return xz, nil
}

func xzReadCmd(ctx context.Context, r io.Reader) (*exec.Cmd, io.Reader, error) {
	cmd := exec.CommandContext(ctx, xzPath, "--quiet", "--decompress", "--stdout", "-")
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
