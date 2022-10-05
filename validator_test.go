package validator

import (
	"testing"
	"time"
)

type Person struct {
	Id        int64
	Name      string
	BirthDate time.Time
}
type Address struct {
	Id      int64
	Street  string
	City    string
	Region  string
	Country string
	Email   string
}

func TestIsRequired(t *testing.T) {
	p := Person{Name: ""}
	err := New(p).IsRequired("Name").Errors
	e := "Name is required"
	if len(err) > 0 {
		a := err["Name"][0].Error()
		if a != e {
			t.Errorf("expected %s, got %s", e, a)
		}
	}
}

func TestIsRequiredExcept(t *testing.T) {
	b := time.Now()
	b = b.Add(time.Duration(-20))
	p := Person{Name: "Adams Smith", BirthDate: b}
	err := New(p).IsRequiredExcept("BirthDate").Errors
	e := "Id is required"
	if len(err) > 0 {
		a := err["Id"][0].Error()
		if e != a {
			t.Errorf("expected %s, got %s", e, a)
		}
	}

}
func TestIsEmail(t *testing.T) {
	t.Run("contains invalid email address", func(t *testing.T) {
		p := Address{Email: "smith@"}
		err := New(p).IsEmail("Email").Errors
		e := "Email contains invalid email address"
		a := err["Email"][0].Error()
		if a != e {
			t.Errorf("%s was expected but found %s", e, a)
		}
	})
	t.Run("contains valid email address", func(t *testing.T) {
		p := Address{Email: "smith@andela"}
		err := New(p).IsEmail("Email").Errors
		e := "Email contains invalid email address"
		if len(err) != 0 {
			a := err["Email"][0].Error()
			t.Errorf("%s was expected but found %s", e, a)
		}
	})
}
