package test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/kom0055/shuttle/core"
	"github.com/stretchr/testify/assert"
)

func (a *A1) MarshalJSON() ([]byte, error) {
	return core.MarshalJson(a)
}

func (a *A1) UnmarshalJSON(b []byte) error {
	*a = A1{}
	err := core.UnmarshalJson(a, b)
	if err != nil {
		return err
	}
	return nil
}

func (a *A2) MarshalJSON() ([]byte, error) {
	return core.MarshalJson(a)
}
func (a *A2) UnmarshalJSON(b []byte) error {

	*a = A2{}
	err := core.UnmarshalJson(a, b)
	if err != nil {
		return err
	}
	return nil
}

func TestJsonMarshal(t *testing.T) {

	{
		b, err := json.MarshalIndent(&in1, "", "\t")
		if err != nil {
			t.Fatal(err)
		}
		_ = os.WriteFile("data/in1.json", b, 0666)
	}
	{

		b, err := json.MarshalIndent(&in2, "", "\t")
		if err != nil {
			t.Fatal(err)
		}
		_ = os.WriteFile("data/in2.json", b, 0666)
	}
}

func TestJsonUnmarshal(t *testing.T) {
	{
		b, err := os.ReadFile("data/in1.json")
		if err != nil {
			t.Fatal(err)
		}
		i1 := A1{}
		err = json.Unmarshal(b, &i1)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(i1)
		assert.Equal(t, in1, i1)

	}
	{
		b, err := os.ReadFile("data/in2.json")
		if err != nil {
			t.Fatal(err)
		}
		i2 := A2{}
		err = json.Unmarshal(b, &i2)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, in2, i2)
		t.Log(i2)
	}
}
