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
