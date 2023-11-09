package core

import (
	"encoding/json"
)

type JsonWrapper struct {
	Kind  string          `json:"kind,omitempty" yaml:"kind,omitempty"`
	Value json.RawMessage `json:"value,omitempty" yaml:"value,omitempty"`
}

func (c *Wrapper) UnmarshalJSON(b []byte) error {
	jsonWrapper := new(JsonWrapper)
	return c.unmarshal(func() (string, error) {
		if err := json.Unmarshal(b, jsonWrapper); err != nil {
			return "", err
		}
		return jsonWrapper.Kind, nil
	}, func(a any) error {
		return json.Unmarshal(jsonWrapper.Value, a)
	})
}
