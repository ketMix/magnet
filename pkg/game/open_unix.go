//go:build aix || darwin || dragonfly || freebsd || (js && wasm) || linux || netbsd || openbsd || solaris

package game

import (
	"os/exec"
)

func OpenFile(path string) error {
	program := "xdg-open"
	// First check if "xdg-open" is available.
	_, err := exec.LookPath("xdg-open")
	// Otherwise default to "open".
	if err != nil {
		program = "open"
	}

	cmd := exec.Command(program, path)
	err = cmd.Start()
	if err != nil {
		return err
	}
	return nil
}
