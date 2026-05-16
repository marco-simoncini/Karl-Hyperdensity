package primitives

import (
	"fmt"
	"strconv"
	"strings"
)

func valueAt(obj map[string]interface{}, path ...string) (interface{}, bool) {
	if len(path) == 0 || obj == nil {
		return nil, false
	}
	var current interface{} = obj
	for _, key := range path {
		m, ok := current.(map[string]interface{})
		if !ok || m == nil {
			return nil, false
		}
		child, exists := m[key]
		if !exists {
			return nil, false
		}
		current = child
	}
	return current, true
}

// StringAt returns the string at path or false if missing/invalid.
func StringAt(obj map[string]interface{}, path ...string) (string, bool) {
	v, ok := valueAt(obj, path...)
	if !ok {
		return "", false
	}
	switch t := v.(type) {
	case string:
		return t, true
	case fmt.Stringer:
		return t.String(), true
	default:
		return fmt.Sprint(v), true
	}
}

// Int64At returns int64 at path (coerces numeric types) or false.
func Int64At(obj map[string]interface{}, path ...string) (int64, bool) {
	v, ok := valueAt(obj, path...)
	if !ok {
		return 0, false
	}
	return coerceInt64(v)
}

// Float64At returns float64 at path or false.
func Float64At(obj map[string]interface{}, path ...string) (float64, bool) {
	v, ok := valueAt(obj, path...)
	if !ok {
		return 0, false
	}
	switch t := v.(type) {
	case float64:
		return t, true
	case float32:
		return float64(t), true
	case int:
		return float64(t), true
	case int64:
		return float64(t), true
	case int32:
		return float64(t), true
	case string:
		f, err := strconv.ParseFloat(strings.TrimSpace(t), 64)
		if err != nil {
			return 0, false
		}
		return f, true
	default:
		return 0, false
	}
}

// MapAt returns nested map at path or false.
func MapAt(obj map[string]interface{}, path ...string) (map[string]interface{}, bool) {
	v, ok := valueAt(obj, path...)
	if !ok {
		return nil, false
	}
	m, ok := v.(map[string]interface{})
	if !ok || m == nil {
		return nil, false
	}
	return m, true
}

// SliceAt returns slice at path or false.
func SliceAt(obj map[string]interface{}, path ...string) ([]interface{}, bool) {
	v, ok := valueAt(obj, path...)
	if !ok {
		return nil, false
	}
	s, ok := v.([]interface{})
	if !ok {
		return nil, false
	}
	return s, true
}

func coerceInt64(v interface{}) (int64, bool) {
	switch t := v.(type) {
	case int:
		return int64(t), true
	case int64:
		return t, true
	case int32:
		return int64(t), true
	case float64:
		return int64(t), true
	case float32:
		return int64(t), true
	case string:
		i, err := strconv.ParseInt(strings.TrimSpace(t), 10, 64)
		if err != nil {
			return 0, false
		}
		return i, true
	default:
		return 0, false
	}
}
