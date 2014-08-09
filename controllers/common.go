package controllers

import (
	"encoding/gob"
	"regexp"
	"strconv"

	"code.google.com/p/go.crypto/bcrypt"
	"github.com/martini-contrib/binding"
)

func init() {
	gob.Register(binding.Errors{})
	gob.Register(map[string]interface{}{})
}

func generatePassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), 10)
}

func formatErr(errors interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	switch errors := errors.(type) {
	case binding.Errors:
		for _, err := range errors {
			for _, field := range err.Fields() {
				if _, ok := result[field]; !ok {
					result[field] = err.Error()
				}
			}
		}
	case map[string]interface{}:
		result = errors
	}

	return result
}

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

func (v *Validator) err(msg string) {
	v.validation.Errors.Add([]string{v.key}, "", msg)
}

// MinLength checks minimum length of a field
func (v *Validator) MinLength(length int, msg string) *Validator {
	if v.len() < length {
		if msg == "" {
			msg = "Minimum length is " + strconv.Itoa(length)
		}

		v.err(msg)
	}

	return v
}

// MaxLength checks maximum length of a field
func (v *Validator) MaxLength(length int, msg string) *Validator {
	if v.len() > length {
		if msg == "" {
			msg = "Maximum length is " + strconv.Itoa(length)
		}

		v.err(msg)
	}
	return v
}

// Length checks whether a field's length is between a specified range
func (v *Validator) Length(min int, max int, msg string) *Validator {
	if len := v.len(); len > max || len < min {
		if msg == "" {
			msg = "Length is between " + strconv.Itoa(min) + "~" + strconv.Itoa(max)
		}

		v.err(msg)
	}

	return v
}

// Email checks whether a field is a valid email
func (v *Validator) Email(msg string) *Validator {
	if !rEmail.MatchString(v.str()) {
		if msg == "" {
			msg = "Email is invalid"
		}

		v.err(msg)
	}

	return v
}

// Equal checks whether a field is equal to a specified value
func (v *Validator) Equal(value interface{}, msg string) *Validator {
	equality := false

	switch x := value.(type) {
	case string:
		equality = v.str() == x
	case *string:
		equality = v.str() == *x
	}

	if !equality {
		if msg == "" {
			msg = v.key + " does not match"
		}

		v.err(msg)
	}

	return v
}
