package main

import "os/exec"
import "strconv"

// ping returns nil if the specified IP address
// is responsive to pings, and returns an error
// otherwise.
func ping(ipAddr string, timeout int) (err error) {
	cmd := exec.Command("ping", pingCountFlag, "1", "-w", strconv.Itoa(timeout), ipAddr)
	if err := cmd.Start(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}
