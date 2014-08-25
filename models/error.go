package models

const (
	UnknownError    = "UnknownError"
	EmailTakenError = "EmailTakenError"
)

type ModelError struct {
	code string
}

func (e ModelError) Error() string {
	return e.code
}
