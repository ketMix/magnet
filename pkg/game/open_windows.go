package game

import (
	"os/exec"
	"path/filepath"
)

func OpenFile(path string) error {
	// This seems bad.
	cmd := exec.Command("rundll32.exe", "url.dll,FileProtocolHandler", path)
	err := cmd.Start()
	if err != nil {
		return err
	}
	return nil
}
