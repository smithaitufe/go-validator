package validator

import (
	"reflect"
	"strings"
)

func getStruct(model interface{}) reflect.Value {
	var ptr reflect.Value
	var v reflect.Value
	v = reflect.ValueOf(model)
	if v.Type().Kind() == reflect.Ptr {
		ptr = v
		v = ptr.Elem()
	} else {
		ptr = reflect.New(reflect.TypeOf(model))
		temp := ptr.Elem()
		temp.Set(v)
	}
	return v
}

func isComplexNumber(f reflect.Value) bool {
	switch f.Kind() {
	case reflect.Complex64, reflect.Complex128:
		return true
	}
	return false
}

func isWholeNumber(f reflect.Value) bool {
	switch f.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	}
	return false
}
func isFloatingNumber(f reflect.Value) bool {
	switch f.Kind() {
	case reflect.Float32, reflect.Complex64, reflect.Complex128:
		return true
	}
	return false
}

func isNumber(f reflect.Value) bool {
	return isWholeNumber(f) || isFloatingNumber(f) || isComplexNumber(f)
}
func isBlank(field reflect.Value) bool {
	switch field.Kind() {
	case reflect.String:
		return len(strings.TrimSpace(field.String())) <= 0
	case reflect.Chan, reflect.Map, reflect.Slice, reflect.Array:
		return field.Len() <= 0
	case reflect.Ptr, reflect.Interface, reflect.Func:
		return field.IsNil()
	default:
		return field.IsValid() && field.Interface() == reflect.Zero(field.Type()).Interface()
	}
}

func addError(v *Validator, field string, e error) {
	if errs, ok := v.Errors[field]; ok {
		v.Errors[field] = append(errs, e)
	} else {
		v.Errors[field] = []error{e}
	}
}
