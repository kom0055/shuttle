package core

import (
	"fmt"
	"reflect"
	"strings"
)

var (
	emptyValue       = reflect.Value{}
	WrapperType      = reflect.TypeOf(Wrapper{})
	WrapperSliceType = reflect.TypeOf([]Wrapper{})
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

func wrapValue(oriVal reflect.Value) (reflect.Value, error) {

	val := oriVal
	typ := RevealInterface(val).Type()
	name, err := globalProvider.getNameByType(typ)
	if err != nil {
		return emptyValue, err
	}

	wrapperPtr := reflect.New(WrapperType)
	wrapperVal := wrapperPtr.Elem()
	wrapperVal.Field(0).SetString(name)
	wrapperVal.Field(1).Set(val)
	return wrapperVal, nil
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

func shuttleMarshal(in any) (any, error) {
	inVal := reflect.ValueOf(in)
	inVal = RevealValue(inVal)
	inType := inVal.Type()
	outType := globalProvider.getMarshalType(inType)
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
				wrappedVal, err := wrapValue(oriVal.Index(i))
				if err != nil {
					return nil, err
				}

				newSliceVal.Index(i).Set(wrappedVal)
			}
			outVal.Field(i).Set(newSliceVal)
			continue
		}

		wrappedVal, err := wrapValue(oriVal)
		if err != nil {
			return nil, err
		}
		outVal.Field(i).Set(wrappedVal)
		continue
	}
	outIf := outVal.Interface()
	return outIf, nil
}

func shuttleUnmarshal(out any, unmarshal func(any, reflect.Type, reflect.Type) error) error {
	outVal := reflect.ValueOf(out)
	if outVal.Kind() != reflect.Ptr {
		return fmt.Errorf("discovery: can only unmarshal into a struct pointer: %T", out)
	}
	outVal = RevealValue(outVal)
	if outVal.Kind() != reflect.Struct {
		return fmt.Errorf("discovery: can only unmarshal into a struct pointer: %T", out)
	}
	outTyp := outVal.Type()
	inType := globalProvider.getMarshalType(outTyp)
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
