package xz

import (
	"bytes"
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
