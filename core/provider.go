package core

import (
	"encoding/json"
	"fmt"
	"reflect"

	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/vmihailenco/msgpack/v5"
)

const (
	shuttleTag   = "shuttle"
	tagFieldWrap = "wrap"
)

var (
	globalProvider = newProvider()
)

func newProvider() *provider {
	return &provider{
		type2NameMap:   cmap.NewStringer[reflect.Type, string](),
		name2TypeMap:   cmap.New[reflect.Type](),
		marshalTypeMap: cmap.NewStringer[reflect.Type, reflect.Type](),
	}
}

type provider struct {
	type2NameMap cmap.ConcurrentMap[reflect.Type, string]
	name2TypeMap cmap.ConcurrentMap[string, reflect.Type]

	marshalTypeMap cmap.ConcurrentMap[reflect.Type, reflect.Type]
}

func (p *provider) getMarshalType(inType reflect.Type) reflect.Type {
	marshalType, ok := p.marshalTypeMap.Get(inType)
	if ok {
		return marshalType
	}

	return p.marshalTypeMap.Upsert(inType, nil, func(exist bool, valueInMap reflect.Type, newValue reflect.Type) reflect.Type {
		if exist {
			return valueInMap
		}
		return createMarshalType(inType)
	})

}

func (p *provider) getTypeByName(string2 string) (reflect.Type, error) {
	t, ok := p.name2TypeMap.Get(string2)
	if !ok {
		return nil, fmt.Errorf("type %s not found", string2)
	}
	return t, nil
}

func (p *provider) getNameByType(t reflect.Type) (string, error) {
	name, ok := p.type2NameMap.Get(t)
	if !ok {
		return "", fmt.Errorf("type %s not found", t)
	}
	return name, nil
}

func (p *provider) setBidNameType(name string, t reflect.Type) error {
	if ok := p.type2NameMap.SetIfAbsent(t, name); !ok {
		return fmt.Errorf("type %s already exist", t)
	}
	if ok := p.name2TypeMap.SetIfAbsent(name, t); !ok {
		return fmt.Errorf("name %s already exist", name)
	}
	return nil
}

func (p *provider) RegisterType(name string, id any) error {
	val := reflect.ValueOf(id)
	typ := val.Type()
	if err := p.setBidNameType(name, typ); err != nil {
		return err
	}
	return nil
}

func (p *provider) MarshalJson(in any) ([]byte, error) {
	outIf, err := p.shuttleMarshal(in)
	if err != nil {
		return nil, err
	}
	return json.Marshal(outIf)
}

func (p *provider) UnmarshalJson(out any, raw json.RawMessage) error {
	return p.shuttleUnmarshal(out, func(a any, inType, outTyp reflect.Type) error {
		return json.Unmarshal(raw, a)
	})

}

func (p *provider) MarshalMsgPack(in any) ([]byte, error) {
	outIf, err := p.shuttleMarshal(in)
	if err != nil {
		return nil, err
	}
	return msgpack.Marshal(outIf)
}

func (p *provider) UnmarshalMsgPack(out any, b []byte) error {
	return p.shuttleUnmarshal(out, func(a any, inType, outTyp reflect.Type) error {
		return msgpack.Unmarshal(b, a)
	})

}

func (p *provider) MarshalYaml(in any) (any, error) {
	return p.shuttleMarshal(in)
}

func (p *provider) UnmarshalYaml(out any, unmarshal func(any) error) error {

	return p.shuttleUnmarshal(out, func(a any, inType, outTyp reflect.Type) error {
		if err := unmarshal(a); err != nil {
			return ReplaceYAMLTypeError(err, inType, outTyp)
		}
		return nil
	})

}

func (p *provider) shuttleMarshal(in any) (any, error) {
	inVal := reflect.ValueOf(in)
	inVal = RevealValue(inVal)
	inType := inVal.Type()
	outType := p.getMarshalType(inType)
	outPtr := reflect.New(outType)
	outVal := outPtr.Elem()
	for i, n := 0, inType.NumField(); i < n; i++ {
		inField := inType.Field(i)
		if !inField.IsExported() {
			continue
		}
		oriVal := inVal.Field(i)
		if !needWrap(inField) {
			outVal.Field(i).Set(oriVal)
			continue
		}

		fieldType := oriVal.Type()
		if fieldType.Kind() == reflect.Slice {
			newSliceVal := reflect.MakeSlice(WrapperSliceType, oriVal.Len(), oriVal.Len())
			for i, n := 0, oriVal.Len(); i < n; i++ {
				wrappedVal, err := p.wrapValue(oriVal.Index(i))
				if err != nil {
					return nil, err
				}

				newSliceVal.Index(i).Set(wrappedVal)
			}
			outVal.Field(i).Set(newSliceVal)
			continue
		}

		wrappedVal, err := p.wrapValue(oriVal)
		if err != nil {
			return nil, err
		}
		outVal.Field(i).Set(wrappedVal)
		continue
	}
	outIf := outVal.Interface()
	return outIf, nil
}

func (p *provider) shuttleUnmarshal(out any, unmarshal func(any, reflect.Type, reflect.Type) error) error {
	outVal := reflect.ValueOf(out)
	if outVal.Kind() != reflect.Ptr {
		return fmt.Errorf("discovery: can only unmarshal into a struct pointer: %T", out)
	}
	outVal = RevealValue(outVal)
	if outVal.Kind() != reflect.Struct {
		return fmt.Errorf("discovery: can only unmarshal into a struct pointer: %T", out)
	}
	outTyp := outVal.Type()
	inType := p.getMarshalType(outTyp)
	inPtr := reflect.New(inType)
	inVal := inPtr.Elem()

	if err := unmarshal(inPtr.Interface(), inType, outTyp); err != nil {
		return err
	}

	for i, n := 0, inType.NumField(); i < n; i++ {
		inField := inType.Field(i)
		if !inField.IsExported() {
			continue
		}

		inValIdxi := inVal.Field(i)
		fieldTypeIdxi := inValIdxi.Type()
		fieldTypeIdxi = RevealType(fieldTypeIdxi)
		if fieldTypeIdxi == WrapperSliceType {
			outValIdxi := outVal.Field(i)
			outValTypeIdxi := outValIdxi.Type()

			cvTyp := fieldTypeIdxi.Elem()
			if cvTyp != WrapperType {
				outVal.Field(i).Set(inValIdxi)
				continue
			}
			if inValIdxi.Len() == 0 {
				continue
			}
			newSliceVal := reflect.MakeSlice(outValTypeIdxi, inValIdxi.Len(), inValIdxi.Len())
			for i, n := 0, inValIdxi.Len(); i < n; i++ {

				val := inValIdxi.Index(i)
				field1 := val.Field(1)
				field1 = RevealInterface(field1)
				newSliceVal.Index(i).Set(field1)
			}
			outVal.Field(i).Set(newSliceVal)
			continue
		}

		if fieldTypeIdxi == WrapperType {
			val := inValIdxi

			field1 := val.Field(1)
			field1 = RevealInterface(field1)
			outVal.Field(i).Set(field1)

			continue
		}
		outVal.Field(i).Set(inValIdxi)
	}
	return nil
}

func (p *provider) wrapValue(oriVal reflect.Value) (reflect.Value, error) {

	val := oriVal
	typ := RevealInterface(val).Type()
	name, err := p.getNameByType(typ)
	if err != nil {
		return emptyValue, err
	}

	wrapper := Wrapper{
		Kind:  name,
		Value: oriVal.Interface(),
	}
	return reflect.ValueOf(wrapper), nil
}
