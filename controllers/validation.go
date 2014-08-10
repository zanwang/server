package controllers

import (
	"regexp"
	"strconv"

	"github.com/martini-contrib/binding"
)

var (
	rEmail = regexp.MustCompile(".+@.+\\..+")
)

// Validation stores the pointer of errors
type Validation struct {
	Errors *binding.Errors
}

// Validate creates a new validator
func (v *Validation) Validate(field interface{}, key string) *Validator {
	return &Validator{key: key, value: field, validation: v}
}

// Validator validates a single field
type Validator struct {
	key        string
	value      interface{}
	validation *Validation
}

func (v *Validator) len() int {
	switch x := v.value.(type) {
	case string:
		return len(x)
	case *string:
		return len(*x)
	case []interface{}:
		return len(x)
	case *[]interface{}:
		return len(*x)
	}

	return 0
}

func (v *Validator) str() string {
	switch x := v.value.(type) {
	case string:
		return x
	case *string:
		return *x
	}

	return ""
}

func (v *Validator) err(code string, msg string) {
	v.validation.Errors.Add([]string{v.key}, code, msg)
}

// Required checks whether a field is not empty
func (v *Validator) Required(msg string) *Validator {
	if v.str() == "" {
		if msg == "" {
			msg = "Required"
		}

		v.err("110", msg)
	}

	return v
}

// MinLength checks minimum length of a field
func (v *Validator) MinLength(length int, msg string) *Validator {
	if v.str() != "" && v.len() < length {
		if msg == "" {
			msg = "Minimum length is " + strconv.Itoa(length)
		}

		v.err("111", msg)
	}

	return v
}

// MaxLength checks maximum length of a field
func (v *Validator) MaxLength(length int, msg string) *Validator {
	if v.str() != "" && v.len() > length {
		if msg == "" {
			msg = "Maximum length is " + strconv.Itoa(length)
		}

		v.err("112", msg)
	}
	return v
}

// Length checks whether a field's length is between a specified range
func (v *Validator) Length(min int, max int, msg string) *Validator {
	if len := v.len(); v.str() != "" && (len > max || len < min) {
		if msg == "" {
			msg = "Length isn't between " + strconv.Itoa(min) + " and " + strconv.Itoa(max)
		}

		v.err("113", msg)
	}

	return v
}

// Email checks whether a field is a valid email
func (v *Validator) Email(msg string) *Validator {
	if v.str() != "" && !rEmail.MatchString(v.str()) {
		if msg == "" {
			msg = "Email is invalid"
		}

		v.err("114", msg)
	}

	return v
}
