package test

import (
	"os"
	"testing"

	"github.com/kom0055/shuttle/core"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func (a *A1) MarshalYAML() (any, error) {
	return core.MarshalYaml(a)
}

func (a *A1) UnmarshalYAML(unmarshal func(any) error) error {
	*a = A1{}
	err := core.UnmarshalYaml(a, unmarshal)
	if err != nil {
		return err
	}
	return nil
}

func (a *A2) MarshalYAML() (any, error) {
	return core.MarshalYaml(a)
}

func (a *A2) UnmarshalYAML(unmarshal func(any) error) error {
	*a = A2{}
	err := core.UnmarshalYaml(a, unmarshal)
	if err != nil {
		return err
	}
	return nil
}

func TestYamlMarshal(t *testing.T) {

	{
		b, err := yaml.Marshal(&in1)
		if err != nil {
			t.Fatal(err)
		}
		_ = os.WriteFile("data/in1.yaml", b, 0666)
	}
	{
		b, err := yaml.Marshal(&in2)
		if err != nil {
			t.Fatal(err)
		}
		_ = os.WriteFile("data/in2.yaml", b, 0666)
	}
}

func TestYamlUnmarshal(t *testing.T) {
	{
		b, err := os.ReadFile("data/in1.yaml")
		if err != nil {
			t.Fatal(err)
		}
		i1 := A1{}
		err = yaml.Unmarshal(b, &i1)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(i1)
		assert.Equal(t, in1, i1)

	}
	{
		b, err := os.ReadFile("data/in2.yaml")
		if err != nil {
			t.Fatal(err)
		}
		i2 := A2{}
		err = yaml.Unmarshal(b, &i2)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, in2, i2)
		t.Log(i2)
	}
}
