package types

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestSerializeAndDeserializeYaml(t *testing.T) {
	original := map[string]interface{}{
		"john": "doe",
		"nested": map[string]interface{}{
			"foo": "bar",
		},
	}

	var yamlOrig YAML = original

	serialized, err := yamlOrig.Value()
	assert.NilError(t, err)

	var yamlNew *YAML = &YAML{}

	assert.NilError(t, yamlNew.Scan([]byte(serialized.(string))))
	assert.DeepEqual(t, &yamlOrig, yamlNew)
}

func TestYamlDeserializationError(t *testing.T) {
	brokenYaml := `foo: bar
faa faz
kaput: true
`

	var yaml *YAML = &YAML{}

	err := yaml.Scan([]byte(brokenYaml))
	assert.Assert(t, err != nil)

	println(err.Error())
}
