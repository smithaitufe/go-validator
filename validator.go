package validator

import (
	"fmt"
	"net/mail"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Validator struct {
	Errors map[string][]error
	model  reflect.Value
}

func New(model interface{}) *Validator {
	m := getStruct(model)
	return &Validator{model: m, Errors: make(map[string][]error)}
}
func (v *Validator) Error() string {
	errs := ""
	for k, es := range v.Errors {
		s := []string{}
		for _, err := range es {
			s = append(s, err.Error())
		}
		errs += fmt.Sprintf("%s: %s", k, strings.Join(s, ","))
	}
	return errs
}
func (v *Validator) IsValid() bool {
	return len(v.Errors) == 0
}
func (v *Validator) IsRequired(fields ...string) *Validator {
	for _, field := range fields {
		f := v.model.FieldByName(field)
		if isBlank(f) {
			err := fmt.Errorf("%s is required", field)
			if errs, ok := v.Errors[field]; ok {
				v.Errors[field] = append(errs, err)
			} else {
				v.Errors[field] = []error{err}
			}
		}
	}
	return v
}
func (v *Validator) IsRequiredExcept(fields ...string) *Validator {
	for i := 0; i < v.model.NumField(); i++ {
		f := v.model.Field(i)
		field := v.model.Type().Field(i).Name
		if strings.Contains(strings.Join(fields, ""), field) {
			continue
		}
		if isBlank(f) {
			err := fmt.Errorf("%s is required", field)
			addError(v, field, err)
		}
	}
	return v
}
func (v *Validator) IsEmail(fields ...string) *Validator {
	for _, field := range fields {
		f := v.model.FieldByName(field)
		if !isBlank(f) {
			_, e := mail.ParseAddress(f.String())
			if e != nil {
				err := fmt.Errorf("%s contains invalid email address", field)
				addError(v, field, err)
			}
		}
	}
	return v
}
func (v *Validator) IsEqualOrBelowMax(field string, max interface{}) *Validator {
	f := v.model.FieldByName(field)
	if !isBlank(f) && isNumber(f) {
		if i, err := strconv.ParseFloat(f.String(), 64); err != nil {
			addError(v, field, err)
		} else {
			if i > max.(float64) {
				addError(v, field, fmt.Errorf("%s cannot exceed %v", field, max))
			}
		}
	}
	return v
}
func (v *Validator) IsEqualOrAboveMin(field string, min interface{}) *Validator {
	f := v.model.FieldByName(field)
	if !isBlank(f) && isNumber(f) {
		if i, err := strconv.ParseFloat(f.String(), 64); err != nil {
			addError(v, field, err)
		} else {
			if i < min.(float64) {
				addError(v, field, fmt.Errorf("%s cannot exceed %v", field, min))
			}
		}
	}
	return v
}
func (v *Validator) IsWholeNumber(fields ...string) *Validator {
	for _, field := range fields {
		f := v.model.FieldByName(field)
		if !isBlank(f) && !isWholeNumber(f) {
			addError(v, field, fmt.Errorf("%s must be a whole number", field))
		}
	}
	return v
}
func (v *Validator) IsFloatingNumber(field string) *Validator {
	f := v.model.FieldByName(field)
	if !isBlank(f) && !isFloatingNumber(f) {
		addError(v, field, fmt.Errorf("%s must be a floating point number", field))
	}
	return v
}
func (v *Validator) IsNumber(field string) *Validator {
	f := v.model.FieldByName(field)
	if !isBlank(f) && !isNumber(f) {
		addError(v, field, fmt.Errorf("%s must be a number", field))
	}
	return v
}

func (v *Validator) IsDate(fields ...string) *Validator {
	for _, field := range fields {
		f := v.model.FieldByName(field)
		if !isBlank(f) {
			l := "2006-01-02"
			_, err := time.Parse(l, f.String())
			if err != nil {
				addError(v, field, err)
			}

		}
	}
	return v
}
