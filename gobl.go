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

	Task("watch-only").
		Watch("cmd/*/*", "cmd/*/*/*", "pkg/*/*", "pkg/*/*/*", "pkg/*/*/*/*").
		Run("build")

	Task("host").
		Watch("./magnet" + exe).
		Signaler(SigQuit).
		Run("runServer")

	Task("join").
		Watch("./magnet" + exe).
		Signaler(SigQuit).
		Run("runClient")

	Task("runServer").
		Exec("./magnet"+exe, "-n", "server", "--host", ":9999")

	Task("runClient").
		Sleep("500ms").
		Exec("./magnet"+exe, "-n", "client", "--join", "localhost:9999")

	Task("run").
		Exec("./magnet" + exe)

	Go()
}
