// +build windows

package ping

import (
	"os/exec"
	"strconv"
)

func command(ipAddr string, timeout int) *exec.Cmd {
	return exec.Command("ping", "-n", "1", "-w", strconv.Itoa(timeout*1000), ipAddr)
}
