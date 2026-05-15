module github.com/marco-simoncini/Karl-Hyperdensity

go 1.22

require (
	github.com/marco-simoncini/Karl-Hyperdensity/pkg/hyperdensity/contractkit v0.0.0
	gopkg.in/yaml.v3 v3.0.1
)

replace github.com/marco-simoncini/Karl-Hyperdensity/pkg/hyperdensity/contractkit => ./pkg/hyperdensity/contractkit
