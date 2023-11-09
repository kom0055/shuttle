package core

import (
	"fmt"
	"reflect"
)

type KindGetter func() (string, error)
type ValueGetter func(any) error

type Wrapper struct {
	Kind  string `json:"kind,omitempty" yaml:"kind,omitempty"`
	Value any    `json:"value,omitempty" yaml:"value,omitempty"`
}

func (c *Wrapper) unmarshal(kindGetter KindGetter, valueGetter ValueGetter) error {
	if kindGetter == nil || valueGetter == nil {
		return fmt.Errorf("kindGetter or valueGetter is nil")
	}
	kind, err := kindGetter()
	if err != nil {
		return err
	}
	outPutType, err := globalProvider.getTypeByName(kind)
	if err != nil {
		return err
	}
	outPutPtr := reflect.New(outPutType)
	outPutKind := outPutType.Kind()
	switch outPutKind {
	case reflect.Ptr:
		elemTyp := outPutType.Elem()
		outPutPtr.Elem().Set(reflect.New(elemTyp))
	case reflect.Struct:
	default:
		return fmt.Errorf("unsupport reflect kind: %v, type: %v", outPutKind.String(), outPutType.String())
	}
	if err := valueGetter(outPutPtr.Interface()); err != nil {
		return err
	}
	v := outPutPtr.Elem().Interface().(any)
	*c = Wrapper{}
	c.Value = v
	c.Kind = kind
	return nil
}
