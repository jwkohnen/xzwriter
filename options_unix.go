//go:build linux || darwin
// +build linux darwin

package xzwriter

// WithSeparateProcessGroup set's the process group of the `xz` subprocess to its own, separate process group.  When
// the program using this library is started in a shell session, hitting CTRL+C will send an interrupt signal to both
// processes.  That means that the `xz` process terminates immediately without reading STDIN to its end, instead
// writes to the xzwriter fails with ERRPIPE.
//
// If the program using this library wants to handle the SIGINT gracefully, one needs to prevent the shell from sending
// the SIGINT to the xz subprocess also.  Running `xz` in a separate process group achieves that.
func WithSeparateProcessGroup() Option {
	return func(xz *XZWriter) error {
		xz.opts.separateProcessGroup = true

		return nil
	}
}
