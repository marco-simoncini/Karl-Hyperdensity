package primitives

import (
	"encoding/json"
	"strings"
)

const PrimitivesPackageVersion = "v0.0.0-sprint49"

// ContractSample is one input/output row in the primitives golden contract.
type ContractSample struct {
	Category string      `json:"category"`
	Input    interface{} `json:"input,omitempty"`
	Path     []string    `json:"path,omitempty"`
	Output   interface{} `json:"output,omitempty"`
	OK       bool        `json:"ok"`
}

// ContractDocument is the deterministic Sprint 49 primitives contract snapshot.
type ContractDocument struct {
	ContractVersion string           `json:"contractVersion"`
	NestedSamples   []ContractSample `json:"nestedSamples"`
	QuantitySamples []ContractSample `json:"quantitySamples"`
}

// DefaultContractDocument returns canonical contract samples for golden tests.
func DefaultContractDocument() ContractDocument {
	obj := map[string]interface{}{
		"metadata": map[string]interface{}{
			"name": "demo",
			"uid":  "abc-123",
		},
		"spec": map[string]interface{}{
			"replicas": int64(3),
			"containers": []interface{}{
				map[string]interface{}{"name": "app", "image": "demo:v1"},
			},
		},
	}
	doc := ContractDocument{
		ContractVersion: PrimitivesPackageVersion,
		NestedSamples: []ContractSample{
			{Category: "nested_string_found", Path: []string{"metadata", "name"}, Input: obj, OK: true},
			{Category: "nested_string_missing", Path: []string{"metadata", "missing"}, Input: obj, OK: false},
			{Category: "nested_int64_found", Path: []string{"spec", "replicas"}, Input: obj, OK: true},
			{Category: "nested_slice_found", Path: []string{"spec", "containers"}, Input: obj, OK: true},
			{Category: "nested_map_found", Path: []string{"metadata"}, Input: obj, OK: true},
		},
		QuantitySamples: []ContractSample{
			{Category: "cpu_100m", Input: "100m", OK: true},
			{Category: "cpu_1", Input: "1", OK: true},
			{Category: "cpu_2", Input: "2", OK: true},
			{Category: "memory_128Mi", Input: "128Mi", OK: true},
			{Category: "memory_1Gi", Input: "1Gi", OK: true},
			{Category: "memory_512Ki", Input: "512Ki", OK: true},
			{Category: "memory_plain_bytes", Input: "1000", OK: true},
			{Category: "cpu_invalid", Input: "not-cpu", OK: false},
			{Category: "memory_invalid", Input: "not-memory", OK: false},
		},
	}
	fillNestedOutputs(&doc, obj)
	fillQuantityOutputs(&doc)
	return doc
}

func fillNestedOutputs(doc *ContractDocument, obj map[string]interface{}) {
	for i := range doc.NestedSamples {
		s := &doc.NestedSamples[i]
		switch s.Category {
		case "nested_string_found":
			v, ok := StringAt(obj, s.Path...)
			s.OK = ok
			if ok {
				s.Output = v
			}
		case "nested_string_missing":
			_, ok := StringAt(obj, s.Path...)
			s.OK = ok
		case "nested_int64_found":
			v, ok := Int64At(obj, s.Path...)
			s.OK = ok
			if ok {
				s.Output = v
			}
		case "nested_slice_found":
			v, ok := SliceAt(obj, s.Path...)
			s.OK = ok
			if ok {
				s.Output = len(v)
			}
		case "nested_map_found":
			_, ok := MapAt(obj, s.Path...)
			s.OK = ok
			if ok {
				s.Output = "map"
			}
		}
	}
}

func fillQuantityOutputs(doc *ContractDocument) {
	for i := range doc.QuantitySamples {
		s := &doc.QuantitySamples[i]
		in, _ := s.Input.(string)
		switch {
		case strings.HasPrefix(s.Category, "cpu_"):
			q, mc, ok := NormalizeCPUQuantity(in)
			s.OK = ok
			if ok {
				s.Output = map[string]interface{}{"quantity": q, "millicores": mc}
			}
		case strings.HasPrefix(s.Category, "memory_"):
			q, b, ok := NormalizeMemoryQuantity(in)
			s.OK = ok
			if ok {
				s.Output = map[string]interface{}{"quantity": q, "bytes": b}
			}
		}
	}
}

// CanonicalContractJSON returns marshaled contract with computed outputs.
func CanonicalContractJSON() ([]byte, error) {
	return json.Marshal(DefaultContractDocument())
}
