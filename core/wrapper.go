package core

import (
	"fmt"
	"reflect"
)

type _KindGetter func() (string, error)
type _ValueGetter func(any) error

type _Wrapper struct {
	Kind        string `json:"kind,omitempty" yaml:"kind,omitempty"`
	Value       any    `json:"value,omitempty" yaml:"value,omitempty"`
	valueGetter _ValueGetter
}

func (c *_Wrapper) unmarshal(kindGetter _KindGetter, valueGetter _ValueGetter) error {
	if kindGetter == nil || valueGetter == nil {
		return fmt.Errorf("kindGetter or valueGetter is nil")
	}
	kind, err := kindGetter()
	if err != nil {
		return err
	}
	*c = _Wrapper{}
	c.Kind = kind
	c.valueGetter = valueGetter
	return nil
}

func (c *_Wrapper) getValue(p *provider) (any, error) {
	outPutType, err := p.getTypeByName(c.Kind)
	if err != nil {
		return nil, err
	}
	outPutPtr := reflect.New(outPutType)
	outPutKind := outPutType.Kind()
	switch outPutKind {
	case reflect.Ptr:
		elemTyp := outPutType.Elem()
		outPutPtr.Elem().Set(reflect.New(elemTyp))
	case reflect.Struct:
	default:
		return nil, fmt.Errorf("unsupport reflect kind: %v, type: %v", outPutKind.String(), outPutType.String())
	}
	if err := c.valueGetter(outPutPtr.Interface()); err != nil {
		return nil, err
	}
	v := outPutPtr.Elem().Interface()
	return v, nil
}
