package main

import "testing"

func TestPing(t *testing.T) {

	if err := isPingable("localhost"); err != nil {
		t.Error("\"localhost\" should be pingable: " + err.Error())
	}

	if err := isPingable("tsohlacol"); err == nil {
		t.Error("It would be very surprising if \"tsohlacol\" was pingable")
	}
}
