package controllers

import (
	"net"
	"regexp"
	"strconv"
	"strings"

	"github.com/martini-contrib/binding"
)

var (
	rEmail      = regexp.MustCompile(".+@.+\\..+")
	rDomain     = regexp.MustCompile("\\.[a-zA-Z]{2,}$")
	rDomainName = regexp.MustCompile("^[a-zA-Z]+[a-zA-Z\\d\\-]*$")
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
	case int:
		return strconv.Itoa(x)
	case *int:
		return strconv.Itoa(*x)
	case int64:
		return strconv.FormatInt(x, 10)
	case *int64:
		return strconv.FormatInt(*x, 10)
	}

	return ""
}

func (v *Validator) int() int {
	var num int
	var err error

	switch x := v.value.(type) {
	case string:
		num, err = strconv.Atoi(x)
	case *string:
		num, err = strconv.Atoi(*x)
	case int:
		num = x
	case *int:
		num = *x
	case int64:
		num = int(x)
	case *int64:
		num = int(*x)
	}

	if err != nil {
		panic(err)
	}

	return num
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
			msg = "Length must be between " + strconv.Itoa(min) + " and " + strconv.Itoa(max)
		}

		v.err("113", msg)
	}

	return v
}

// Email checks whether a field is a valid email
func (v *Validator) Email(msg string) *Validator {
	if str := v.str(); str != "" && !rEmail.MatchString(str) {
		if msg == "" {
			msg = "Email is invalid"
		}

		v.err("114", msg)
	}

	return v
}

func (v *Validator) Min(num int, msg string) *Validator {
	if v.int() > num {
		if msg == "" {
			msg = "Minimum value is " + strconv.Itoa(num)
		}

		v.err("115", msg)
	}

	return v
}

func (v *Validator) Max(num int, msg string) *Validator {
	if v.int() > num {
		if msg == "" {
			msg = "Maximum value is " + strconv.Itoa(num)
		}

		v.err("116", msg)
	}

	return v
}

func (v *Validator) Within(min int, max int, msg string) *Validator {
	num := v.int()

	if num > max && num < min {
		if msg == "" {
			msg = "Value must be between " + strconv.Itoa(min) + "~" + strconv.Itoa(max)
		}

		v.err("117", msg)
	}

	return v
}

func (v *Validator) Without(min int, max int, msg string) *Validator {
	num := v.int()

	if num <= max && num >= min {
		if msg == "" {
			msg = "Value must be not between " + strconv.Itoa(min) + "~" + strconv.Itoa(max)
		}

		v.err("118", msg)
	}

	return v
}

func (v *Validator) IsIn(arr []string, msg string) *Validator {
	var inarr bool
	str := v.str()

	for _, x := range arr {
		if x == str {
			inarr = true
			break
		}
	}

	if !inarr {
		if msg == "" {
			msg = "Value must be one of [" + strings.Join(arr, ", ") + "]"
		}

		v.err("119", msg)
	}

	return v
}

func (v *Validator) NotIn(arr []string, msg string) *Validator {
	str := v.str()

	for _, x := range arr {
		if x == str {
			if msg == "" {
				msg = "Value must not be one of [" + strings.Join(arr, ", ") + "]"
			}

			v.err("120", msg)
			break
		}
	}

	return v
}

func (v *Validator) IP(msg string) *Validator {
	if str := v.str(); str != "" {
		if ip := net.ParseIP(str); ip == nil {
			if msg == "" {
				msg = "IP is invalid"
			}

			v.err("121", msg)
		}
	}

	return v
}

func (v *Validator) Domain(msg string) *Validator {
	if str := v.str(); str != "" && !rDomain.MatchString(str) {
		if msg == "" {
			msg = "Domain is invalid"
		}

		v.err("123", msg)
	}

	return v
}

func (v *Validator) DomainName(msg string) *Validator {
	if str := v.str(); str != "" && !rDomainName.MatchString(str) {
		if msg == "" {
			msg = "Domain name is invalid"
		}

		v.err("124", msg)
	}

	return v
}
