package core

import (
	"encoding/json"
	"reflect"
)

func MarshalJson(in any) ([]byte, error) {
	outIf, err := shuttleMarshal(in)
	if err != nil {
		return nil, err
	}
	return json.Marshal(outIf)
}

func UnmarshalJson(out any, raw json.RawMessage) error {

	return shuttleUnmarshal(out, func(a any, inType, outTyp reflect.Type) error {
		return json.Unmarshal(raw, a)
	})

}
