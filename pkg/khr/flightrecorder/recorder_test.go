package flightrecorder

import "testing"

func TestFlightRecorder(t *testing.T) {
	Reset()
	Record("register", "host ok", "")
	if len(Snapshot()) != 1 {
		t.Fatal("expected one event")
	}
}
