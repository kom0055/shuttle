package test

import (
	"os"
	"testing"

	"github.com/kom0055/shuttle/core"
	"github.com/stretchr/testify/assert"
	"github.com/vmihailenco/msgpack/v5"
)

func (a *A1) MarshalMsgpack() ([]byte, error) {
	return core.MarshalMsgPack(a)
}

func (a *A1) UnmarshalMsgpack(b []byte) error {
	*a = A1{}
	err := core.UnmarshalMsgPack(a, b)
	if err != nil {
		return err
	}
	return nil
}

func (a *A2) MarshalMsgpack() ([]byte, error) {
	return core.MarshalMsgPack(a)
}
func (a *A2) UnmarshalMsgpack(b []byte) error {

	*a = A2{}
	err := core.UnmarshalMsgPack(a, b)
	if err != nil {
		return err
	}
	return nil
}

func TestMsgPackMarshal(t *testing.T) {

	{
		b, err := msgpack.Marshal(&in1)
		if err != nil {
			t.Fatal(err)
		}
		_ = os.WriteFile("data/in1.msgpack", b, 0666)
	}
	{

		b, err := msgpack.Marshal(&in2)
		if err != nil {
			t.Fatal(err)
		}
		_ = os.WriteFile("data/in2.msgpack", b, 0666)
	}
}

func TestMsgPackUnmarshal(t *testing.T) {
	{
		b, err := os.ReadFile("data/in1.msgpack")
		if err != nil {
			t.Fatal(err)
		}
		i1 := A1{}
		err = msgpack.Unmarshal(b, &i1)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(i1)
		assert.Equal(t, in1, i1)

	}
	{
		b, err := os.ReadFile("data/in2.msgpack")
		if err != nil {
			t.Fatal(err)
		}
		i2 := A2{}
		err = msgpack.Unmarshal(b, &i2)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, in2, i2)
		t.Log(i2)
	}
}
