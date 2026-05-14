package main

import (
	"encoding/json"
	"os"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/realmarkettick"
)

func main() {
	doc, err := realmarkettick.ReferenceSurface()
	if err != nil {
		panic(err)
	}
	raw, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile("examples/runtime-market-tick-real-shell-inputs-no-apply-reference.json", raw, 0644); err != nil {
		panic(err)
	}
}
