package core

import (
	"errors"

	"gopkg.in/yaml.v3"
)

var (
	invalidYamlErr = errors.New("invalid yaml")
)

func (c *Wrapper) UnmarshalYAML(value *yaml.Node) error {

	return c.unmarshal(func() (string, error) {
		if value == nil {
			return "", nil
		}

		if len(value.Content) != 4 {
			return "", invalidYamlErr
		}
		return value.Content[1].Value, nil
	}, func(a any) error {
		return value.Content[3].Decode(a)
	})
}
