package forms

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Form creates a custom form struct, embeds an url.Values object
type Form struct {
	url.Values
	Errors errors
}

// New initializes a form struct
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Valid returns true if there are no errors, otherwise false
func (f *Form) IsValid() bool {
	return len(f.Errors) == 0
}

// Required checks if field is not empty
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.AddError(field, "This field cannot be blank")
		}
	}
}

// MinLength checks for minimal lenght of characters in field
func (f *Form) MinLength(field string, length int, req *http.Request) bool {
	fieldVal := req.Form.Get(field)

	if len(fieldVal) < length {
		f.Errors.AddError(field, fmt.Sprintf("Minimal characters length is %d", length))
		return false
	}

	return true
}

//Has checks if form field is in request and not empty
func (f *Form) Has(field string, req *http.Request) bool {
	fieldVal := req.Form.Get(field)

	return fieldVal != ""
}

//IsEmail checks if e-mail is in valid format
func (f *Form) IsEmail(field string, req *http.Request) {
	fieldVal := req.Form.Get(field)
	validate := validator.New()

	err := validate.Var(fieldVal, "email")
	if err != nil {
		f.Errors.AddError(field, "Email must be a valid email address")
	}
}
