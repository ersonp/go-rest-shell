//go:build windows
// +build windows

package shell

import "os/exec"

func Execute(command string) *exec.Cmd {
	return exec.Command("powershell", "-Command", command)
}
