package main

import (
	"os/exec"
)

func OpenBrowser(url string) bool {
	args := []string{"cmd", "/c", "start"}

	cmd := exec.Command(args[0], append(args[1:], url)...)

	return cmd.Start() == nil
}
