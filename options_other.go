//go:build !linux && !darwin
// +build !linux,!darwin

package xzwriter

import "syscall"

func sysProcAttr() *syscall.SysProcAttr { return nil}
