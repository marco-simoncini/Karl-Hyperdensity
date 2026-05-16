package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/shellcontinuity"
)

func main() {
	ctx := flag.String("cluster-context", "", "kubectl context")
	ns := flag.String("namespace", "khr-runtime-sandbox", "pod namespace")
	selector := flag.String("selector", "app=khr-native-live-target", "pod label selector")
	out := flag.String("out", "", "write snapshot JSON")
	flag.Parse()
	if *ctx == "" {
		fmt.Fprintln(os.Stderr, "cluster-context required")
		os.Exit(2)
	}
	pod, uid, cid, err := podIdentity(*ctx, *ns, *selector)
	if err != nil {
		fatal(err)
	}
	snap := shellcontinuity.SnapshotFromWorkload(*ns, pod, uid, cid)
	b, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		fatal(err)
	}
	b = append(b, '\n')
	if *out == "" {
		os.Stdout.Write(b)
		return
	}
	if err := os.WriteFile(*out, b, 0o644); err != nil {
		fatal(err)
	}
}

func podIdentity(ctx, ns, selector string) (pod, uid, cid string, err error) {
	pod, err = kubectl(ctx, ns, "jsonpath={.items[0].metadata.name}", selector)
	if err != nil {
		return
	}
	uid, err = kubectl(ctx, ns, "jsonpath={.items[0].metadata.uid}", selector)
	if err != nil {
		return
	}
	cid, err = kubectl(ctx, ns, "jsonpath={.items[0].status.containerStatuses[0].containerID}", selector)
	return
}

func kubectl(ctx, ns, jp, selector string) (string, error) {
	cmd := exec.Command("kubectl", "--context", ctx, "-n", ns,
		"get", "pods", "-l", selector, "-o", jp)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
