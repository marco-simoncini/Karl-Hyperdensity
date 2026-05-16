package primitives

import (
	"strconv"
	"strings"
)

// NormalizeCPUQuantity parses minimal CPU quantity forms (stdlib-only contract).
// Supports: "100m", whole cores "1", "2". Unknown input returns ok=false.
func NormalizeCPUQuantity(input string) (quantity string, millicores int64, ok bool) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return "", 0, false
	}
	if strings.HasSuffix(trimmed, "m") {
		num := strings.TrimSpace(strings.TrimSuffix(trimmed, "m"))
		if num == "" {
			return "", 0, false
		}
		v, err := strconv.ParseInt(num, 10, 64)
		if err != nil || v < 0 {
			return "", 0, false
		}
		return trimmed, v, true
	}
	v, err := strconv.ParseInt(trimmed, 10, 64)
	if err != nil || v < 0 {
		return "", 0, false
	}
	return trimmed, v * 1000, true
}

// NormalizeMemoryQuantity parses minimal memory forms (stdlib-only contract).
// Supports: "128Mi", "1Gi", "512Ki", plain bytes "1000". Unknown returns ok=false.
func NormalizeMemoryQuantity(input string) (quantity string, bytes int64, ok bool) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return "", 0, false
	}
	type suffix struct {
		suffix string
		mult   int64
	}
	for _, s := range []suffix{
		{"Gi", 1024 * 1024 * 1024},
		{"Mi", 1024 * 1024},
		{"Ki", 1024},
	} {
		if strings.HasSuffix(trimmed, s.suffix) {
			num := strings.TrimSpace(strings.TrimSuffix(trimmed, s.suffix))
			if num == "" {
				return "", 0, false
			}
			v, err := strconv.ParseInt(num, 10, 64)
			if err != nil || v < 0 {
				return "", 0, false
			}
			return trimmed, v * s.mult, true
		}
	}
	v, err := strconv.ParseInt(trimmed, 10, 64)
	if err != nil || v < 0 {
		return "", 0, false
	}
	return trimmed, v, true
}
