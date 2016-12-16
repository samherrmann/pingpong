package network_test

import (
	"testing"

	"github.com/samherrmann/pingpong/network"
)

func TestPing(t *testing.T) {

	if err := network.Ping("localhost", 1); err != nil {
		t.Error("\"localhost\" should be pingable: " + err.Error())
	}

	if err := network.Ping("tsohlacol", 1); err == nil {
		t.Error("It would be very surprising if \"tsohlacol\" was pingable")
	}
}
