# shuttle
Shuttle Marshal/Unmarshal

It's hard to deserialize bytes to a specified derived struct implements some interfaces.

So you have to write a lot of code to do this.

Referenced and enhanced the idea from [prometheus](https://github.com/prometheus/prometheus).

Weave like shuttle.

It supports `json` ,`msgpack` and `yaml` now, but generated idl like `protobuf`, `thrift` is not supported.


## Usage

You could see the test cases in `test/` folder.

```Golang

type Flyable interface{
	Fly()
}

type Bird struct{
    Name string `json:"name"`
}

func (a *Bird) Fly() {

}

type Plane struct{
	Num string `json:"num"`
}

func (a *Plane) Fly() {

}


// step 1. register your structs 

core.RegisterType("Bird", &Bird{})

core.RegisterType("Plane", &Plane{})

// step 2. modify your struct field tags, add shuttle:",wrap", 
// and Implement Marshaler and Unmarshaler of json or yaml

type SomeStructA struct {
    F1 Flyable `json:"F1" shuttle:",wrap"`
    F2 Flyable `json:"F2", shuttle:",wrap""`
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

// step3. serialize your struct SomeStructA to bytes, and deserialize bytes to SomeStructA

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

// {"F1":{"kind":"Bird","value":{"name":"Flying Bird"}},"F2":{"kind":"Plane","value":{"num":"Plane 01"}}}

```

