package core

import (
	"encoding/json"
)

func MarshalJson(in any) ([]byte, error) {
	return globalProvider.MarshalJson(in)
}

func UnmarshalJson(out any, raw json.RawMessage) error {
	return globalProvider.UnmarshalJson(out, raw)
}
