package main

import (
	"encoding/json"
	"os"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/canarycohort"
)

func main() {
	doc, err := canarycohort.ReferenceSurface()
	if err != nil {
		panic(err)
	}
	raw, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile("examples/production-canary-cohort-expansion-reference.json", raw, 0644); err != nil {
		panic(err)
	}
}
