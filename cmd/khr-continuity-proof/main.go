package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/khr/shellcontinuity"
)

func main() {
	beforePath := flag.String("before", "", "continuity snapshot JSON (before)")
	afterPath := flag.String("after", "", "continuity snapshot JSON (after)")
	outPath := flag.String("out", "", "write proof JSON")
	flag.Parse()
	if *beforePath == "" || *afterPath == "" {
		fmt.Fprintln(os.Stderr, "usage: khr-continuity-proof -before file -after file [-out file]")
		os.Exit(2)
	}
	before, err := loadSnapshot(*beforePath)
	if err != nil {
		fatal(err)
	}
	after, err := loadSnapshot(*afterPath)
	if err != nil {
		fatal(err)
	}
	proof := shellcontinuity.Compare(before, after)
	if *outPath == "" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(proof)
		return
	}
	b, _ := json.MarshalIndent(proof, "", "  ")
	b = append(b, '\n')
	if err := os.WriteFile(*outPath, b, 0o644); err != nil {
		fatal(err)
	}
}

func loadSnapshot(path string) (shellcontinuity.Snapshot, error) {
	var s shellcontinuity.Snapshot
	b, err := os.ReadFile(path)
	if err != nil {
		return s, err
	}
	return s, json.Unmarshal(b, &s)
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
