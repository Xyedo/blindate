package validator

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/invopop/validation"
)

func UniqueByStructFields(visitor func(i int) any) func(value any) error {
	return func(value any) error {
		value, isNil := validation.Indirect(value)
		if isNil || validation.IsEmpty(value) {
			return nil
		}

		fields := reflect.ValueOf(value)
		if !(fields.Kind() == reflect.Slice || fields.Kind() == reflect.Array) {
			panic(fmt.Sprintf("invalid field type %T", fields.Interface()))
		}

		vs := make([]any, 0, fields.Len())

		for i := 0; i < fields.Len(); i++ {
			vs = append(vs, visitor(i))
		}

		return unique(vs)
	}
}
func Unique(value any) error {
	values, isNil := validation.Indirect(value)
	if isNil || validation.IsEmpty(values) {
		return nil
	}

	return unique(value)
}

func unique(value any) error {
	field := reflect.ValueOf(value)
	v := reflect.ValueOf(struct{}{})

	errValue := errors.New("must be unique")

	switch field.Kind() {
	case reflect.Slice, reflect.Array:
		elem := field.Type().Elem()
		if elem.Kind() == reflect.Pointer {
			elem = elem.Elem()
		}

		m := reflect.MakeMap(reflect.MapOf(elem, v.Type()))

		for i := 0; i < field.Len(); i++ {
			m.SetMapIndex(reflect.Indirect(field.Index(i)), v)
		}

		if field.Len() != m.Len() {
			return errValue
		}
	case reflect.Map:
		var m reflect.Value

		if field.Type().Elem().Kind() == reflect.Pointer {
			m = reflect.MakeMap(reflect.MapOf(field.Type().Elem().Elem(), v.Type()))
		} else {
			m = reflect.MakeMap(reflect.MapOf(field.Type().Elem(), v.Type()))
		}

		for _, k := range field.MapKeys() {
			m.SetMapIndex(reflect.Indirect(field.MapIndex(k)), v)
		}
		if field.Len() != m.Len() {
			return errValue
		}
	default:
		panic(fmt.Sprintf("invalid field type %T", field.Interface()))
	}

	return nil
}
