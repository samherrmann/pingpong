package ping

import "errors"

// Run pings the specified address and returns nil
// if a response is received. An error is returned
// if no response is received.
func Run(addr string, timeout int) (err error) {
	cmd := command(addr, timeout)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(string(out))
	}
	return nil
}
