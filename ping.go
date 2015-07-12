package main

import "os/exec"

// ping returns nil if the specified IP address
// is responsive to pings, and returns an error
// otherwise.
func ping(ipAddr string) (err error) {
	cmd := exec.Command("ping", pingCountFlag, "1", "-w", "1", ipAddr)
	if err := cmd.Start(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}
