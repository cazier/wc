package models

import (
	"reflect"
)

// I'm not really sure this works reilably

func IsNil(m any) bool {
	switch reflect.TypeOf(m).Kind() {
	case reflect.Map:
		return nilMap(m)
	case reflect.Slice:
		return nilSlice(m)
	case reflect.Struct:
		return nilStruct(m)
	case reflect.String:
		return m == ""
	default:
		return reflect.ValueOf(m).IsZero() || reflect.ValueOf(m).IsNil()
	}
}

func nilMap(m any) bool {
	v := reflect.ValueOf(m)
	if v.Len() == 0 {
		return true
	}

	iter := v.MapRange()
	for iter.Next() {
		v := iter.Value().Interface()
		return IsNil(v)
	}
	return true
}

func nilSlice(m any) bool {
	v := reflect.ValueOf(m)
	if v.Len() == 0 {
		return true
	}

	for index := 0; index < v.Len(); index++ {
		if v.Index(index) != reflect.Zero(v.Index(index).Type()).Interface() {
			return false
		}
	}
	return true
}

func nilStruct(m interface{}) bool {
	v := reflect.ValueOf(m)

	for index := 0; index < v.NumField(); index++ {
		field := v.Field(index)

		if field.Interface() != reflect.Zero(field.Type()).Interface() {
			return false
		}
	}

	return true
}
