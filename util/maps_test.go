package util

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestCopyMap(t *testing.T) {
	map1 := map[string]interface{}{
		"john": "doe",
		"piet": "snot",
	}

	map2, err := CopyMap(map1)
	assert.NilError(t, err)

	assert.DeepEqual(t, map1, map2)

	map2["jane"] = "doe"
	_, ok := map1["jane"]
	assert.Assert(t, !ok)

	map2["john"] = "legend"
	assert.Assert(t, map1["john"] != map2["john"])
}

func TestMergeMap_Simple(t *testing.T) {
	m1 := map[string]interface{}{
		"override_me": "m1",
		"in_m1":       "m1",
	}
	m2 := map[string]interface{}{
		"override_me": "m2",
		"in_m2":       "m2",
	}

	out, err := MergeMaps(m1, m2)
	assert.NilError(t, err)
	assert.Equal(t, out["override_me"], "m2")
	assert.Equal(t, out["in_m1"], "m1")
	assert.Equal(t, out["in_m2"], "m2")
}

func TestMergeMap_Nested(t *testing.T) {
	m1 := map[string]interface{}{
		"nested": map[string]interface{}{
			"key1":  "m1",
			"other": "m1",
		},
	}
	m2 := map[string]interface{}{
		"nested": map[string]interface{}{
			"key2":  "m2",
			"other": "m2",
		},
	}

	out, err := MergeMaps(m1, m2)
	assert.NilError(t, err)
	assert.Equal(t, out["nested"].(map[string]interface{})["other"], "m2")
	assert.Equal(t, out["nested"].(map[string]interface{})["key1"], "m1")
	assert.Equal(t, out["nested"].(map[string]interface{})["key2"], "m2")
}
