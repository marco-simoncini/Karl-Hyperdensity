package host

import "testing"

func TestRegisterHost(t *testing.T) {
	cfg := &Config{}
	cfg.Spec.HostID = "host-1"
	cfg.Spec.LinuxOnly = true
	cfg.Spec.SandboxMode = true
	reg := RegisterHost(cfg)
	if reg.HostID != "host-1" || reg.ProductionSafe {
		t.Fatalf("reg=%+v", reg)
	}
	if reg.Tool != "karl-host-runtime" {
		t.Fatal("tool name")
	}
}

func TestValidateConfigRequiresSandbox(t *testing.T) {
	cfg := &Config{}
	cfg.Spec.HostID = "x"
	cfg.Spec.LinuxOnly = true
	if len(ValidateConfig(cfg)) == 0 {
		t.Fatal("expected sandboxMode error")
	}
}
