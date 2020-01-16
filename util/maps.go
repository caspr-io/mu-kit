package util

import (
	"encoding/json"
)

// Map performs a deep copy of the given map m.
func CopyMap(m map[string]interface{}) (map[string]interface{}, error) {
	jsonMap, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	var c map[string]interface{}

	if err := json.Unmarshal(jsonMap, &c); err != nil {
		return nil, err
	}

	return c, nil
}

func MergeMaps(maps ...map[string]interface{}) (map[string]interface{}, error) {
	merged := map[string]interface{}{}

	for _, m := range maps {
		unmarshalled, err := CopyMap(m)
		if err != nil {
			return nil, err
		}

		merged = mergeMapKV(merged, unmarshalled)
	}

	return merged, nil
}

func mergeMapKV(m1, m2 map[string]interface{}) map[string]interface{} {
	for k, v2 := range m2 {
		if v1, ok := m1[k]; ok {
			m1[k] = merge1(v1, v2)
		} else {
			m1[k] = v2
		}
	}

	return m1
}

// merge1 merges 2 interfaces, x2 wins
func merge1(x1, x2 interface{}) interface{} {
	switch x1 := x1.(type) {
	case map[string]interface{}:
		x2, ok := x2.(map[string]interface{})
		if !ok {
			return x1
		}

		return mergeMapKV(x1, x2)
	case nil:
		// merge(nil, map[string]interface{...}) -> map[string]interface{...}
		x2, ok := x2.(map[string]interface{})
		if ok {
			return x2
		}
	}

	return x2
}
