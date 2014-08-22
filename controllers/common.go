package controllers

import (
	"strconv"

	"code.google.com/p/go.crypto/bcrypt"
	"github.com/martini-contrib/binding"
)

func generatePassword(password string) string {
	if hash, err := bcrypt.GenerateFromPassword([]byte(password), 10); err != nil {
		panic(err)
	} else {
		return string(hash)
	}
}

// APIError stores API error
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewErr(fields []string, code string, msg string) binding.Errors {
	errors := binding.Errors{}
	errors.Add(fields, code, msg)

	return errors
}

func FormatErr(errors interface{}) map[string]interface{} {
	result := make(map[string]APIError)

	switch errors := errors.(type) {
	case binding.Errors:
		for _, err := range errors {
			fields := err.Fields()

			if len(fields) == 0 {
				fields = []string{"common"}
			}

			for _, field := range fields {
				if _, ok := result[field]; !ok {
					msg := err.Error()

					switch err.Kind() {
					case binding.TypeError:
						result[field] = APIError{125, msg}
					case binding.ContentTypeError:
						result[field] = APIError{126, msg}
					case binding.DeserializationError:
						result[field] = APIError{127, msg}
					default:
						if code, err := strconv.Atoi(err.Kind()); err == nil {
							result[field] = APIError{code, msg}
						}
					}
				}
			}
		}
	}

	return map[string]interface{}{
		"errors": result,
	}
}

func toint64(str string) int64 {
	if num, err := strconv.ParseInt(str, 10, 64); err != nil {
		panic(err)
	} else {
		return num
	}
}
