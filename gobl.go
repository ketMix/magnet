package main

import (
	"runtime"

	. "github.com/kettek/gobl"
)

func main() {
	var exe string
	if runtime.GOOS == "windows" {
		exe = ".exe"
	}

	Task("build").
		Exec("go", "build", "./cmd/magnet")

	Task("watch").
		Watch("cmd/*/*", "cmd/*/*/*", "pkg/*/*", "pkg/*/*/*", "pkg/*/*/*/*").
		Signaler(SigQuit).
		Run("build").
		Run("run")

	Task("run").
		Exec("./magnet" + exe)

	Go()
}
