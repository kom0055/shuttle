package test

import (
	"testing"

	"github.com/kom0055/shuttle/core"
)

type A interface {
	MustImplementA()
}

type A1 struct {
	Field1 string   `json:"field1,omitempty" yaml:"field1,omitempty"`
	Filed2 []A      `json:"filed2,omitempty" yaml:"filed2,omitempty" shuttle:",wrap"`
	Field3 int      `json:"field3,omitempty" yaml:"field3,omitempty"`
	Field4 A        `json:"field4,omitempty" yaml:"field4,omitempty" shuttle:",wrap"`
	Filed5 []A      `json:"filed5,omitempty" yaml:"filed5,omitempty" shuttle:",wrap"`
	Field6 []string `json:"field6,omitempty" yaml:"field6,omitempty"`
}

type A2 struct {
	Field1 []A `json:"field1,omitempty" yaml:"field1,omitempty" shuttle:",wrap"`
	Field2 A   `json:"field2,omitempty" yaml:"field2,omitempty" shuttle:",wrap"`
}

type impl1 struct {
	What string `json:"what,omitempty" yaml:"what,omitempty"`
}

func (i *impl1) Name() string {
	return "impl1"
}

func (i *impl1) MustImplementA() {

}

type impl2 struct {
	How string `json:"how,omitempty" yaml:"how,omitempty"`
}

func (i *impl2) Name() string {
	return "impl2"
}

func (i *impl2) MustImplementA() {

}

type impl3 struct {
	When int64 `json:"when,omitempty" yaml:"when,omitempty"`
}

func (i *impl3) Name() string {
	return "impl3"
}

func (i *impl3) MustImplementA() {

}

type impl5 struct {
	Who string `json:"who,omitempty" yaml:"who,omitempty"`
}

func (i impl5) Name() string {
	return "impl5"
}

func (i impl5) MustImplementA() {

}

var (
	in1 = A1{
		Field1: "field1",
		Filed2: []A{
			&impl2{
				How: "2",
			},
			&impl1{
				What: "3",
			},
			&impl2{
				How: "4",
			},
			impl5{Who: "kun"},
		},
		Field3: 3,
		Field4: &impl1{
			What: "5",
		},
		Filed5: []A{
			&impl1{
				What: "6",
			},
			&impl1{
				What: "7",
			},
			&impl2{
				How: "8",
			},
			&impl2{
				How: "9",
			},
		},
		Field6: []string{"a", "b", "c"},
	}

	in2 = A2{
		Field1: []A{
			&impl1{
				What: "a2what",
			},
			&impl2{How: "a2how"},
			&impl3{When: 1111111},
		},
		Field2: impl5{Who: "peng"},
	}
)

func TestMain(m *testing.M) {
	i1, i2, i3, i5 := &impl1{}, &impl2{}, &impl3{}, impl5{}

	if err := core.RegisterType(i1.Name(), i1); err != nil {
		panic(err)
	}
	if err := core.RegisterType(i2.Name(), i2); err != nil {
		panic(err)
	}
	if err := core.RegisterType(i3.Name(), i3); err != nil {
		panic(err)
	}
	if err := core.RegisterType(i5.Name(), i5); err != nil {
		panic(err)
	}

	m.Run()
}
