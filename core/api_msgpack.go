package core

import (
	"reflect"

	"github.com/vmihailenco/msgpack/v5"
)

func MarshalMsgPack(in any) ([]byte, error) {
	outIf, err := shuttleMarshal(in)
	if err != nil {
		return nil, err
	}
	return msgpack.Marshal(outIf)
}

func UnmarshalMsgPack(out any, b []byte) error {
	return shuttleUnmarshal(out, func(a any, inType, outTyp reflect.Type) error {
		return msgpack.Unmarshal(b, a)
	})

}
