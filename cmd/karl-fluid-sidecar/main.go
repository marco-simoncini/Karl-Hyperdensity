package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/windowsfluidvirt/sidecar"
)

func main() {
	var socketPath string
	var qemuPID string
	var sidecarVersion string
	flag.StringVar(&socketPath, "qmp-socket", "/var/run/kubevirt-private/qemu.sock", "QMP socket path")
	flag.StringVar(&qemuPID, "qemu-pid", "", "QEMU PID hint")
	flag.StringVar(&sidecarVersion, "sidecar-version", "v0-readonly", "sidecar version")
	flag.Parse()

	transport := sidecar.NewSocketTransport(socketPath)
	executor := sidecar.NewReadOnlyExecutor(transport)
	evidence := executor.DiscoverEvidence(socketPath, qemuPID, sidecarVersion)

	encoded, err := json.MarshalIndent(evidence, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to encode qmp evidence: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(encoded))
}
