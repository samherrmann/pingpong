package network

import "errors"

// Ping pings the specified network address.
func Ping(addr string, timeout int) error {
	cmd := ping(addr, timeout)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(string(out))
	}
	return nil
}
