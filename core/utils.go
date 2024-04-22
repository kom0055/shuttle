package core

import (
	"reflect"
	"strings"
)

var (
	emptyValue       = reflect.Value{}
	WrapperType      = reflect.TypeOf(_Wrapper{})
	WrapperSliceType = reflect.TypeOf([]_Wrapper{})
)

func RevealInterface(val reflect.Value) reflect.Value {
	for val.Kind() == reflect.Interface {
		val = val.Elem()
	}
	return val
}

func RevealValue(val reflect.Value) reflect.Value {
	for val.Kind() == reflect.Ptr || val.Kind() == reflect.Interface {
		val = val.Elem()
	}
	return val
}

func RevealType(typ reflect.Type) reflect.Type {
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	return typ
}

func createMarshalType(inTyp reflect.Type) reflect.Type {
	inTyp = RevealType(inTyp)
	var fields []reflect.StructField
	for i, n := 0, inTyp.NumField(); i < n; i++ {
		field := inTyp.Field(i)

		if !needWrap(field) {
			fields = append(fields, field)
			continue
		}
		fieldType := field.Type
		if fieldType.Kind() == reflect.Slice {
			fields = append(fields, reflect.StructField{
				Name: field.Name,
				Type: WrapperSliceType,
				Tag:  field.Tag,
			})
			continue
		}

		fields = append(fields, reflect.StructField{
			Name: field.Name,
			Type: WrapperType,
			Tag:  field.Tag,
		})
	}
	return reflect.StructOf(fields)
}

func needWrap(field reflect.StructField) bool {
	tag := field.Tag.Get(shuttleTag)

	if tag == "" && strings.Index(string(field.Tag), ":") < 0 {
		tag = string(field.Tag)
	}

	if tag == "-" {
		return false
	}

	fields := strings.Split(tag, ",")
	if len(fields) > 1 {
		for _, flag := range fields[1:] {
			switch flag {
			case tagFieldWrap:
				return true
			}
		}
	}

	return false

}
