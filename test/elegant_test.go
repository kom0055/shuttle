package test

import (
	"encoding/json"
	"testing"

	"github.com/kom0055/shuttle/core"
	"github.com/stretchr/testify/assert"
)

type Flyable interface {
	Fly()
}

type Bird struct {
	Name string `json:"name"`
}

func (a *Bird) Fly() {

}

type Plane struct {
	Num string `json:"num"`
}

func (a *Plane) Fly() {

}

type SomeStructA struct {
	F1 Flyable `json:"F1" shuttle:",wrap"`
	F2 Flyable `json:"F2" shuttle:",wrap"`
}

func (a SomeStructA) MarshalJSON() ([]byte, error) {
	return core.MarshalJson(a)
}
func (a *SomeStructA) UnmarshalJSON(b []byte) error {

	*a = SomeStructA{}
	err := core.UnmarshalJson(a, b)
	if err != nil {
		return err
	}
	return nil
}

func TestSerialize(t *testing.T) {
	core.RegisterType("Bird", &Bird{})

	core.RegisterType("Plane", &Plane{})
	core.RegisterType("SomeStructA", &SomeStructA{})

	a1 := SomeStructA{
		F1: &Bird{Name: "Flying Bird"},
		F2: &Plane{Num: "Plane 01"},
	}
	b, err := json.Marshal(a1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(b))
	a2 := SomeStructA{}
	if err := json.Unmarshal(b, &a2); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, a1, a2)
}
