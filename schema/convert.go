package schema

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// Unmarshal YAML to map[string]any instead of map[any]any.
func Unmarshal(in []byte, out any) error {
	var res any

	if err := yaml.Unmarshal(in, &res); err != nil {
		return err
	}
	*out.(*any) = mapValue(res)

	return nil
}

func mapSlice(in []any) []any {
	res := make([]any, len(in))
	for i, v := range in {
		res[i] = mapValue(v)
	}
	return res
}

func mapMap(in map[any]any) map[string]any {
	res := make(map[string]any)
	for k, v := range in {
		res[fmt.Sprintf("%v", k)] = mapValue(v)
	}
	return res
}

func mapValue(v any) any {
	switch v := v.(type) {
	case []any:
		return mapSlice(v)
	case map[any]any:
		return mapMap(v)
	default:
		return v
	}
}
