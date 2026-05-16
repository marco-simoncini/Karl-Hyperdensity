package flightrecorder

import "testing"

func TestFlightRecorder(t *testing.T) {
	Reset()
	InitContext(SessionContext{
		RuntimeSessionID:      "khr-session-test",
		HostRuntimeInstanceID: "khr-inst-test",
		CorrelationID:         "khr-corr-test",
	})
	Record("register", "host ok", "")
	if len(Snapshot()) != 1 {
		t.Fatal("expected one event")
	}
	if Snapshot()[0].RuntimeSessionID != "khr-session-test" {
		t.Fatalf("event=%+v", Snapshot()[0])
	}
}
