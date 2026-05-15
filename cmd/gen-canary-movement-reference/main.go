package main

import (
	"encoding/json"
	"os"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/canarymovement"
)

func main() {
	doc, err := canarymovement.ReferenceSurface()
	if err != nil {
		panic(err)
	}
	raw, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile("examples/production-canary-movement-expansion-reference.json", raw, 0644); err != nil {
		panic(err)
	}
}
