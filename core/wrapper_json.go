package core

import (
	"encoding/json"
)

type _JsonWrapper struct {
	Kind  string          `json:"kind,omitempty" yaml:"kind,omitempty"`
	Value json.RawMessage `json:"value,omitempty" yaml:"value,omitempty"`
}

func (c *_Wrapper) UnmarshalJSON(b []byte) error {
	jsonWrapper := new(_JsonWrapper)
	return c.unmarshal(func() (string, error) {
		if err := json.Unmarshal(b, jsonWrapper); err != nil {
			return "", err
		}
		return jsonWrapper.Kind, nil
	}, func(a any) error {
		return json.Unmarshal(jsonWrapper.Value, a)
	})
}
