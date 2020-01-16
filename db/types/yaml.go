package types

import (
	"database/sql/driver"

	"gopkg.in/yaml.v3"
)

type YAML map[string]interface{}

func (m *YAML) Scan(b interface{}) error {
	if b == nil {
		*m = nil
		return nil
	}

	var t map[string]interface{}

	err := yaml.Unmarshal(b.([]byte), &t)
	if err != nil {
		return err
	}

	*m = t

	return nil
}

func (m YAML) Value() (driver.Value, error) {
	b, err := yaml.Marshal(m)
	if err != nil {
		return nil, err
	}

	return string(b), nil
}
