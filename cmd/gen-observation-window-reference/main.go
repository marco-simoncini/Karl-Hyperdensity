package main

import (
	"encoding/json"
	"os"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/observationwindow"
)

func main() {
	doc, err := observationwindow.ReferenceSurface()
	if err != nil {
		panic(err)
	}
	raw, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile("examples/production-observation-window-multitick-compression-evidence-reference.json", raw, 0644); err != nil {
		panic(err)
	}
}
