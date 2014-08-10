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
	Code    string `json:"code"`
	Message string `json:"message"`
}

func newErr(fields []string, code string, msg string) binding.Errors {
	errors := binding.Errors{}
	errors.Add(fields, code, msg)

	return errors
}

func FormatErr(errors interface{}) map[string]interface{} {
	result := make(map[string][]APIError)

	switch errors := errors.(type) {
	case binding.Errors:
		for _, err := range errors {
			for _, field := range err.Fields() {
				if _, ok := result[field]; !ok {
					result[field] = []APIError{}
				}

				result[field] = append(result[field], APIError{
					Code:    err.Kind(),
					Message: err.Error(),
				})
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
