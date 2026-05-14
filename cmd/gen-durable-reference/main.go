package main

import (
	"os"

	"github.com/marco-simoncini/Karl-Hyperdensity/pkg/marketcontroller"
)

func main() {
	data, err := marketcontroller.MarshalDurableReferenceSurfaceJSON()
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile("examples/durable-controller-state-kubernetes-reconciler-reference.json", data, 0644); err != nil {
		panic(err)
	}
}
