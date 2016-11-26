package ping

// Run pings the specified address and returns nil
// if a response is received. An error is returned
// if no response is received.
func Run(addr string, timeout int) (err error) {
	cmd := command(addr, timeout)
	if err := cmd.Start(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}
