package main

import (
	"os"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/marketcontroller"
)

func main() {
	data, err := marketcontroller.MarshalReferenceSurfaceJSON()
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile("examples/live-controller-reconciliation-execution-loop-reference.json", data, 0644); err != nil {
		panic(err)
	}
}
