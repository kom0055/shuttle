package core

import (
	"fmt"
	"reflect"

	cmap "github.com/orcaman/concurrent-map/v2"
)

const (
	shuttleTag   = "shuttle"
	tagFieldWrap = "wrap"
)

var (
	globalProvider = &provider{
		type2NameMap:   cmap.NewStringer[reflect.Type, string](),
		name2TypeMap:   cmap.New[reflect.Type](),
		marshalTypeMap: cmap.NewStringer[reflect.Type, reflect.Type](),
	}
)

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
