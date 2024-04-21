package core

import (
	"errors"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

func MarshalYaml(in any) (any, error) {
	return globalProvider.MarshalYaml(in)
}

func UnmarshalYaml(out any, unmarshal func(any) error) error {
	return globalProvider.UnmarshalYaml(out, unmarshal)

}

func ReplaceYAMLTypeError(err error, oldTyp, newTyp reflect.Type) error {
	var e *yaml.TypeError
	if errors.As(err, &e) {
		oldStr := oldTyp.String()
		newStr := newTyp.String()
		for i, s := range e.Errors {
			e.Errors[i] = strings.Replace(s, oldStr, newStr, -1)
		}
	}
	return err
}
