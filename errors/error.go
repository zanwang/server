package errors

import "net/http"

const (
	// 1xx: Format error
	Unknown         = 110
	Required        = 111
	ContentType     = 112
	Deserialization = 113
	Type            = 114
	Email           = 120
	URL             = 121
	Alpha           = 122
	Alphanumeric    = 123
	Numeric         = 124
	Hexadecimal     = 125
	HexColor        = 126
	LowerCase       = 127
	UpperCase       = 128
	Int             = 129
	Float           = 130
	Divisble        = 131
	Length          = 132
	MinLength       = 133
	MaxLength       = 134
	UUID            = 135
	CreditCard      = 136
	ISBN            = 137
	JSON            = 138
	Multibyte       = 139
	ASCII           = 140
	FullWidth       = 141
	HalfWidth       = 142
	VariableWidth   = 143
	Base64          = 144
	IP              = 145
	IPv4            = 146
	IPv6            = 147
	MAC             = 148
	Min             = 149
	Max             = 150
	Range           = 151
	DomainName      = 152
	Domain          = 153
	// 2xx: Custom error
	UserNotActivated = 210
	UserActivated    = 211
	EmailUsed        = 212
	DomainUsed       = 213
	WrongPassword    = 214
	PasswordUnset    = 215
	RecordType       = 216
	TokenExpired     = 217
	UserNotExist     = 218
	// Same as http status code
	Unauthorized = 401
	Forbidden    = 403
	NotFound     = 404
	ServerError  = 500
)

type API struct {
	Status  int    `json:"-"`
	Field   string `json:"field,omitempty"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e API) Error() string {
	return e.Message
}

func New(field string, code int, message string) API {
	return API{
		Status:  http.StatusBadRequest,
		Field:   field,
		Code:    code,
		Message: message,
	}
}
