package core

import (
	"github.com/vmihailenco/msgpack/v5"
)

type MsgPackWrapper struct {
	Kind  string             `json:"kind,omitempty" yaml:"kind,omitempty"`
	Value msgpack.RawMessage `json:"value,omitempty" yaml:"value,omitempty"`
}

func (c *Wrapper) UnmarshalMsgpack(b []byte) error {
	msgPackWrapper := new(MsgPackWrapper)
	return c.unmarshal(func() (string, error) {
		if err := msgpack.Unmarshal(b, msgPackWrapper); err != nil {
			return "", err
		}
		return msgPackWrapper.Kind, nil
	}, func(a any) error {
		return msgpack.Unmarshal(msgPackWrapper.Value, a)
	})
}
