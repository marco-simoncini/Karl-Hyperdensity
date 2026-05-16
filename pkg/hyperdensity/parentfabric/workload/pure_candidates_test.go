package workload

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestAppsWorkloadResource(t *testing.T) {
	res, ok := AppsWorkloadResource("Deployment")
	if !ok || res != "deployments" {
		t.Fatalf("Deployment: res=%q ok=%v", res, ok)
	}
	res, ok = AppsWorkloadResource("StatefulSet")
	if !ok || res != "statefulsets" {
		t.Fatalf("StatefulSet: res=%q ok=%v", res, ok)
	}
	res, ok = AppsWorkloadResource("DaemonSet")
	if ok || res != "" {
		t.Fatalf("DaemonSet: res=%q ok=%v", res, ok)
	}
}

func TestPilotWorkloadTerm(t *testing.T) {
	term, ok := PilotWorkloadTerm("Deployment")
	if !ok || term != "Deployment" {
		t.Fatalf("Deployment: term=%q ok=%v", term, ok)
	}
	term, ok = PilotWorkloadTerm("Pod")
	if ok || term != "workload" {
		t.Fatalf("Pod: term=%q ok=%v", term, ok)
	}
}

func TestExecutionSupportsLiveApplyKind(t *testing.T) {
	if !ExecutionSupportsLiveApplyKind("Deployment") {
		t.Fatal("Deployment should be supported")
	}
	if ExecutionSupportsLiveApplyKind("VirtualMachine") {
		t.Fatal("VirtualMachine should not be supported")
	}
}

func TestCanonicalContractJSON_matchesGolden(t *testing.T) {
	got, err := CanonicalContractJSON()
	if err != nil {
		t.Fatal(err)
	}
	want, err := os.ReadFile(filepath.Join("testdata", "workload_pure_candidates_contract.golden.json"))
	if err != nil {
		t.Fatal(err)
	}
	var gotDoc, wantDoc map[string]json.RawMessage
	if err := json.Unmarshal(got, &gotDoc); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(want, &wantDoc); err != nil {
		t.Fatal(err)
	}
	for k, gv := range gotDoc {
		wv, ok := wantDoc[k]
		if !ok {
			t.Fatalf("unexpected contract key %q", k)
		}
		if string(gv) != string(wv) {
			t.Fatalf("contract field %q mismatch:\ngot  %s\nwant %s", k, gv, wv)
		}
	}
}
