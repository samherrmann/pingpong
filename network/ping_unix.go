// +build !windows

package network

import (
	"os/exec"
	"strconv"
)

func ping(addr string, timeout int) *exec.Cmd {
	return exec.Command("ping", "-c", "1", "-w", strconv.Itoa(timeout), addr)
}
