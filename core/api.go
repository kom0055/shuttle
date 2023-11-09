package core

import (
	"reflect"
)

func RegisterType(name string, id any) error {
	val := reflect.ValueOf(id)
	typ := val.Type()
	if err := globalProvider.setBidNameType(name, typ); err != nil {
		return err
	}
	return nil
}
